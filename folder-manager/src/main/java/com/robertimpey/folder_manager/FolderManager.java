package com.robertimpey.folder_manager;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.List;

import jakarta.annotation.Nonnull;

public class FolderManager {
    private final List<String> locations;
    private final List<String> folders;

    public FolderManager(List<String> locations, List<String> folders) {
        this.locations = locations;
        this.folders = folders;
    }

    public static @Nonnull FolderManager create(Path locationsPath, Path foldersPath) throws Exception {
        List<String> locations = readLinesFromPath(locationsPath);
        if (locations.isEmpty()) {
            throw new IllegalArgumentException("Locations cannot be null or empty");
        }

        List<String> folders = readLinesFromPath(foldersPath);
        if (folders.isEmpty()) {
            throw new IllegalArgumentException("Folders cannot be null or empty");
        }

        return new FolderManager(locations, folders);
    }

    public void listManagedFolders(PrintWriter outFile) {
        boolean topOfPrintingLocations = true;
        for (String location : this.locations) {
            var locationPath = Paths.get(location);
            if (Files.exists(locationPath)) {
                boolean topOfPrintingFolders = true;
                for (String folder : this.folders) {

                    var folderPath = locationPath.resolve(folder);
                    if (Files.exists(folderPath)) {
                        if (topOfPrintingLocations) {
                            topOfPrintingLocations = false;
                        } else {
                            if (topOfPrintingFolders) {
                                outFile.println();
                            }
                        }

                        outFile.println(folderPath.toAbsolutePath());
                        topOfPrintingFolders = false;
                    }
                }
            }
        }
    }

    public void generateRobocopyScripts(Path autoGenFolder) throws Exception {
        if (autoGenFolder == null || !Files.exists(autoGenFolder)) {
            throw new IllegalArgumentException("Auto-generated folder does not exist: " + autoGenFolder);
        }

        for (String location1 : this.locations) {
            for (String location2 : this.locations) {
                if (location1.equals(location2)) {
                    continue; // Skip if both locations are the same
                }

                var scriptPath = getScriptPath(autoGenFolder, location1, location2);
                var sourcePath = Paths.get(location1);
                var destinationPath = Paths.get(location2);
                List<String> commonFolders = new ArrayList<>();
                for (String folder : this.folders) {

                    if (Files.exists(sourcePath.resolve(folder)) && Files.exists(destinationPath.resolve(folder))) {
                        commonFolders.add(folder);

                        createRobocopySynchScript(folder, scriptPath.resolve(folder + ".ps1"), sourcePath, destinationPath);
                    }
                }

                if (!commonFolders.isEmpty()) {
                    createAllFoldersRobocopySynchScript(commonFolders, scriptPath.resolve("_all.ps1"), sourcePath, destinationPath);
                }
            }
        }
    }

    private static @Nonnull Path getScriptPath(@Nonnull Path autoGenFolder, @Nonnull String location1, @Nonnull String location2) {
        return autoGenFolder.resolve("synch").resolve(getCleanLocationName(location1)).resolve(getCleanLocationName(location2));
    }

    private static @Nonnull String getCleanLocationName(@Nonnull String location) {
        var allLegal = location.replaceAll("[:\\\\/ ]+", "_");

        return allLegal.replaceAll("_+$", "");
    }

    private static void createRobocopySynchScript(@Nonnull String folder, @Nonnull Path scriptPath, @Nonnull Path sourcePath, @Nonnull Path destinationPath) throws Exception {
        if (!Files.exists(scriptPath.getParent())) {
            Files.createDirectories(scriptPath.getParent());
        }

        System.out.println("Generating " + scriptPath);

        try (PrintWriter writer = new PrintWriter(Files.newBufferedWriter(scriptPath))) {
            writeAutoGenHeader(writer);
            writeSynchScriptFileParams(writer);

            writer.println("Import-Module \"$($env:LOCAL_SCRIPTS)\\_Common\\synch\\Synch.psm1\"");
            writer.println();

            writer.printf("$folder = \"%s\"%n", folder);
            writer.printf("$src = \"%s\"%n", sourcePath.toAbsolutePath());
            writer.printf("$dst = \"%s\"%n", destinationPath.toAbsolutePath());
            writer.println();

            writer.println("Synch $folder $src $dst $logged");
        }
    }

    private static void createAllFoldersRobocopySynchScript(@Nonnull List<String> commonFolders, @Nonnull Path scriptPath, @Nonnull Path sourcePath, @Nonnull Path destinationPath) throws Exception {
        if (!Files.exists(scriptPath.getParent())) {
            Files.createDirectories(scriptPath.getParent());
        }

        System.out.println("Generating " + scriptPath);

        try (PrintWriter writer = new PrintWriter(Files.newBufferedWriter(scriptPath))) {
            writeAutoGenHeader(writer);
            writeSynchScriptFileParams(writer);

            writer.println("Import-Module \"$($env:LOCAL_SCRIPTS)\\_Common\\synch\\Synch.psm1\"");
            writer.println();

            writer.print("$folders = ");

            boolean first = true;
            for (String folder : commonFolders) {
                if (first) {
                    first = false;
                } else {
                    writer.print(", ");
                }
                writer.printf("\"%s\"", folder);
            }
            writer.println();
            writer.println();

            writer.printf("$src = \"%s\"%n", sourcePath.toAbsolutePath());
            writer.printf("$dst = \"%s\"%n", destinationPath.toAbsolutePath());
            writer.println();

            writer.println("foreach ($folder in $folders) {");
            writer.println("    Synch $folder $src $dst $logged");
            writer.println("}");
        }
    }

    private static @Nonnull List<String> readLinesFromPath(@Nonnull Path filePath) throws Exception {
        List<String> lines = new ArrayList<>();
        if (Files.exists(filePath)) {
            try (BufferedReader reader = new BufferedReader(new InputStreamReader(Files.newInputStream(filePath)))) {
                String line;
                while ((line = reader.readLine()) != null) {
                    lines.add(line);
                }
            }
        } else {
            throw new Exception("File does not exist: " + filePath);
        }

        return lines;
    }

    public static void writeAutoGenHeader(@Nonnull PrintWriter outFile) {
        outFile.println("# AUTOGEN'D FILE - DO NOT EDIT");

        var zonedDateTime = LocalDateTime.now().atZone(ZoneId.systemDefault());

        // Define a custom time format
        var formatter = DateTimeFormatter.RFC_1123_DATE_TIME;

        // Format the time
        var formattedDateTime = zonedDateTime.format(formatter);
        outFile.printf("# Created: %s%n%n", formattedDateTime);
    }

    private static void writeSynchScriptFileParams(@Nonnull PrintWriter outFile) {
        outFile.println("param(");
        outFile.println("    [Parameter (Mandatory = $False)]");
        outFile.println("    [switch]$logged = $False");
        outFile.println(")");
        outFile.println();
    }
}
