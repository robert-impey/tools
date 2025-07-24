package com.robertimpey.folder_manager;

import java.util.concurrent.Callable;

import picocli.CommandLine.Command;

@Command(name = "subcommands", subcommands = {ListManagedFoldersCommand.class, GenerateRobocopyScriptsCommand.class}, mixinStandardHelpOptions = true, description = "Subcommands for folder management. Use 'list' to list managed folders or 'robocopy' to generate robocopy scripts.")
public class SubcommandsClasses implements Callable<Integer> {

    @Override
    public Integer call() throws Exception {
        System.out.println("Subcommand needed: 'list' or 'robocopy'");
        return 0;
    }
}
