package com.robertimpey.folder_manager;

import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.concurrent.Callable;

import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

@Command(name = "robocopy", mixinStandardHelpOptions = true, description = "Generates robocopy scripts for managed folders.")
public class GenerateRobocopyScriptsCommand implements Callable<Integer> {
    @Option(names = { "-l", "--locations" }, description = "The locations file")
    private String locationsFile;

    @Option(names = { "-f", "--folders" }, description = "The folders file")
    private String foldersFile;

    @Option(names = { "-a", "--autogen" }, description = "The folder for the autogen'd scripts")
    private String autoGenFolder;

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

        System.out.printf("Reading folders from: %s%n", foldersFile);
        Path foldersPath = Paths.get(foldersFile);

        if (autoGenFolder == null || autoGenFolder.isEmpty()) {
            System.out.println("Auto-generated folder is required. Use -a or --auto-gen-folder to specify it.");
            return 0; // Return 0 for success, but no auto-generated folder to list
        }

        System.out.printf("Reading locations from: %s%n", locationsFile);
        System.out.printf("Reading folders from: %s%n", foldersFile);
        System.out.printf("Using auto-generated folder: %s%n", autoGenFolder);

        FolderManager folderManager = FolderManager.create(locationsPath, foldersPath);
        folderManager.generateRobocopyScripts(Paths.get(autoGenFolder));
        
        return 0;
    }

}
