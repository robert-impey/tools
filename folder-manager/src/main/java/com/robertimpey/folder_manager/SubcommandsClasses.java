package com.robertimpey.folder_manager;

import java.util.concurrent.Callable;

import picocli.CommandLine.Command;

@Command(name = "subcommands", subcommands = { ListManagedFoldersCommand.class })
public class SubcommandsClasses implements Callable<Integer> {

    @Override
    public Integer call() throws Exception {
        System.out.println("Subcommand needed: 'list'");
        return 0;
    }
}
