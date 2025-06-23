package com.robertimpey.folder_manager;

import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.util.concurrent.Callable;

@Command(name = "list", mixinStandardHelpOptions = true, description = "Lists all managed folders in the system.")
public class ListManagedFoldersCommand implements Callable<Integer> {

    @Option(names = { "-l", "--locations" }, description = "The locations file")
    private String locations;

    @Option(names = { "-f", "--folders" }, description = "The folders file")
    private String folders;

    @Option(names = { "-m", "--managed-folders-file" }, description = "The managed folders file")
    private String managedFoldersFile;

    @Override
    public Integer call() throws Exception {
        System.out.printf("Reading locations from: %s%n", locations);   
        System.out.printf("Reading folders from: %s%n", folders);   
        System.out.printf("Writing the list of managed folders to: %s%n", managedFoldersFile);   

        return 0; // Return 0 for success
    }

}
