package com.robertimpey.folder_manager;

import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import picocli.CommandLine;

@SpringBootApplication
public class FolderManagerApplication implements ApplicationRunner {

	public static void main(String[] args) {
		SpringApplication.run(FolderManagerApplication.class, args);
	}

	@Override
	public void run(ApplicationArguments args) throws Exception {
		int exitCode = new CommandLine(new SubcommandsClasses()).execute(args.getSourceArgs());
		if (exitCode != 0) {
			throw new RuntimeException(String.format("Command failed with exit code %d%n", exitCode));
		}
	}
}
