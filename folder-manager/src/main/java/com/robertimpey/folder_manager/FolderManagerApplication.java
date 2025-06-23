package com.robertimpey.folder_manager;

import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class FolderManagerApplication implements ApplicationRunner{

	public static void main(String[] args) {
		SpringApplication.run(FolderManagerApplication.class, args);
	}

	@Override
	public void run(ApplicationArguments args) throws Exception {
		System.out.println("Folder Manager Application has started successfully.");
	}

}
