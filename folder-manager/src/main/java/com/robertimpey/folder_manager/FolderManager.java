package com.robertimpey.folder_manager;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
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
}
