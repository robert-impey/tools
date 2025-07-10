package com.robertimpey.folder_manager;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.List;

public class FolderManager {
    private final List<String> locations;
    private final List<String> folders;

    public FolderManager(List<String> locations, List<String> folders) {
        this.locations = locations;
        this.folders = folders;
    }

    public static FolderManager create(Path locationsPath, Path foldersPath) throws Exception {
        List<String> locations = readLinesFromPath(locationsPath);
        if (locations == null || locations.isEmpty()) {
            throw new IllegalArgumentException("Locations cannot be null or empty");
        }

        List<String> folders = readLinesFromPath(foldersPath);
        if (folders == null || folders.isEmpty()) {
            throw new IllegalArgumentException("Folders cannot be null or empty");
        }

        return new FolderManager(locations, folders);
    }

    public void listManagedFolders(PrintWriter outFile) {
        boolean topOfPrintingLocations = true;
        for (String location : this.locations) {
            Path locationPath = Paths.get(location);
            if (Files.exists(locationPath)) {
                boolean topOfPrintingFolders = true;
                for (String folder : this.folders) {

                    Path folderPath = locationPath.resolve(folder);
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

                for (String folder : this.folders) {
                    Path sourcePath = Paths.get(location1, folder);
                    Path destinationPath = Paths.get(location2, folder);

                    if (Files.exists(sourcePath) && Files.exists(destinationPath)) {
                        // Generate the robocopy script
                        Path scriptPath = getScriptPath(autoGenFolder, location1, location2, folder);
                        createRobocopySynchScript(scriptPath, sourcePath, destinationPath);
                    }
                }
            }
        }
    }

    private static Path getScriptPath(Path autoGenFolder, String location1, String location2, String folder) {
        return autoGenFolder
                .resolve("synch")
                .resolve(getCleanLocationName(location1))
                .resolve(getCleanLocationName(location2))
                .resolve(folder + ".ps1");
    }

    private static String getCleanLocationName(String location) {
        String allLegal = location.replaceAll("[:\\\\\\\\/ ]+", "_");

        String noTrailingUnderscore = allLegal.replaceAll("_+$", "");

        return noTrailingUnderscore;
    }

    private static void createRobocopySynchScript(Path scriptPath, Path sourcePath, Path destinationPath)
            throws Exception {
        if (!Files.exists(scriptPath.getParent())) {
            Files.createDirectories(scriptPath.getParent());
        }

        System.out.println("Generating " + scriptPath);

        try (PrintWriter writer = new PrintWriter(Files.newBufferedWriter(scriptPath))) {
            writeAutoGenHeader(writer);
            writeSynchScriptFileParams(writer);

            writer.println("Import-Module \"$($env:LOCAL_SCRIPTS)\\_Common\\synch\\Synch.psm1\"");
            writer.println();

            writer.printf("$folder = \"%s\"%n", sourcePath.getFileName());
            writer.printf("$src = \"%s\"%n", sourcePath.getParent().toAbsolutePath());
            writer.printf("$dst = \"%s\"%n", destinationPath.getParent().toAbsolutePath());
            writer.println();

            writer.println("Synch $folder $src $dst $logged");
        }
    }

    private static List<String> readLinesFromPath(Path filePath) throws Exception {
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

    public static void writeAutoGenHeader(PrintWriter outFile) {
        outFile.println("# AUTOGEN'D FILE - DO NOT EDIT");

        LocalDateTime localDateTime = LocalDateTime.now();
        ZonedDateTime zonedDateTime = localDateTime.atZone(ZoneId.systemDefault());

        // Define a custom time format
        DateTimeFormatter formatter = DateTimeFormatter.RFC_1123_DATE_TIME;

        // Format the time
        String formattedDateTime = zonedDateTime.format(formatter);
        outFile.printf("# Created: %s%n%n", formattedDateTime);
    }

    private static void writeSynchScriptFileParams(PrintWriter outFile) {
        outFile.println("param(");
        outFile.println("    [Parameter (Mandatory = $False)]");
        outFile.println("    [switch]$logged = $False");
        outFile.println(")");
        outFile.println();
    }
}
