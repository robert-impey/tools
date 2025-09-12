package com.robertimpey.folder_manager;

import java.io.OutputStreamWriter;
import java.io.PrintWriter;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.concurrent.Callable;

import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

@Command(name = "list", mixinStandardHelpOptions = true, description = "Lists all managed folders in the system.")
public class ListManagedFoldersCommand implements Callable<Integer> {

    @Option(names = {"-l", "--locations"}, description = "The locations file")
    private String locationsFile;

    @Option(names = {"-f", "--folders"}, description = "The folders file")
    private String foldersFile;

    @Option(names = {"-m", "--managed-folders-file"}, description = "The managed folders file")
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
        var locationsPath = Paths.get(locationsFile);

        System.out.printf("Reading folders from: %s%n", foldersFile);
        var foldersPath = Paths.get(foldersFile);

        var folderManager = FolderManager.create(locationsPath, foldersPath);

        PrintWriter outFile;
        if (managedFoldersFile == null || managedFoldersFile.isEmpty()) {
            outFile = new PrintWriter(System.out, true);
        } else {
            System.out.printf("Writing the list of managed folders to: %s%n", managedFoldersFile);
            var managedFoldersPath = Paths.get(managedFoldersFile);
            if (!Files.exists(managedFoldersPath)) {
                Files.createDirectories(managedFoldersPath.getParent());
            }
            var outputStream = Files.newOutputStream(managedFoldersPath);
            outFile = new PrintWriter(new OutputStreamWriter(outputStream), true);

            FolderManager.writeAutoGenHeader(outFile);
        }

        folderManager.listManagedFolders(outFile);

        return 0; // Return 0 for success
    }
}
