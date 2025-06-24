package com.robertimpey.folder_manager;

import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Callable;

@Command(name = "list", mixinStandardHelpOptions = true, description = "Lists all managed folders in the system.")
public class ListManagedFoldersCommand implements Callable<Integer> {

    @Option(names = { "-l", "--locations" }, description = "The locations file")
    private String locationsFile;

    @Option(names = { "-f", "--folders" }, description = "The folders file")
    private String foldersFile;

    @Option(names = { "-m", "--managed-folders-file" }, description = "The managed folders file")
    private String managedFoldersFile;

    @Override
    public Integer call() throws Exception {
        if (locationsFile == null || locationsFile.isEmpty()) {
            System.out.println("Locations file is required. Use -l or --locations to specify it.");
            return 0; // Return 0 for success, but no locations to list
        }

        if (foldersFile == null || foldersFile.isEmpty()) {
            System.out.println("Folders file is required. Use -f or --folders to specify it.");
            return 0; // Return 0 for success, but no folders to list
        }

        System.out.printf("Reading locations from: %s%n", locationsFile);
        Path locationsPath = Paths.get(locationsFile);
        List<String> locations = null;

        if (Files.exists(locationsPath)) {
            locations = readLinesFromPath(locationsPath);
        }

        System.out.printf("Reading folders from: %s%n", foldersFile);
        Path foldersPath = Paths.get(foldersFile);
        List<String> folders = null;

        if (Files.exists(foldersPath)) {
            folders = readLinesFromPath(foldersPath);
        }

        if (locations == null || locations.isEmpty() || folders == null || folders.isEmpty()) {
            System.out.println("No locations or folders found. Exiting.");
            return 0; // Return 1 for success, but no managed folders to list
        }

        if (managedFoldersFile == null || managedFoldersFile.isEmpty()) {
            for (String location : locations) {
                for (String folder : folders) {
                    System.out.printf("Managed folder: %s/%s%n", location, folder);
                }
            }
        } else {

            System.out.printf("Writing the list of managed folders to: %s%n", managedFoldersFile);
        }

        return 0; // Return 0 for success
    }

    private List<String> readLinesFromPath(Path filePath) throws Exception {
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
