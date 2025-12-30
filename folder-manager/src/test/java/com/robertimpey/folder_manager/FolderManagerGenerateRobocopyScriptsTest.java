package com.robertimpey.folder_manager;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

public class FolderManagerGenerateRobocopyScriptsTest {
    @TempDir
    Path tempDir;
    
    @Test
    void generateRobocopyScripts_createsExpectedScripts() throws Exception {
        Path loc1 = tempDir.resolve("Loc1");
        Path loc2 = tempDir.resolve("Loc2");
        Files.createDirectories(loc1.resolve("Shared"));
        Files.createDirectories(loc2.resolve("Shared"));
        Path autoGen = tempDir.resolve("AutoGen");
        Files.createDirectories(autoGen);

        FolderManager manager = new FolderManager(
                List.of(loc1.toString(), loc2.toString()),
                List.of("Shared")
        );

        manager.generateRobocopyScripts(autoGen);

        Path scriptDir = FolderManager.getScriptPath(autoGen, loc1.toString(), loc2.toString());
        assertTrue(Files.exists(scriptDir.resolve("Shared.ps1")));
        assertTrue(Files.exists(scriptDir.resolve("_all.ps1")));

        Path reverseScriptDir = FolderManager.getScriptPath(autoGen, loc2.toString(), loc1.toString());
        assertTrue(Files.exists(reverseScriptDir.resolve("Shared.ps1")));
        assertTrue(Files.exists(reverseScriptDir.resolve("_all.ps1")));
    }

    @Test
    void generateRobocopyScripts_throwsIfAutoGenFolderDoesNotExist() {
        Path missing = tempDir.resolve("MissingAutoGen");
        FolderManager manager = new FolderManager(List.of("C:\\Data"), List.of("Folder"));

        IllegalArgumentException ex = assertThrows(
                IllegalArgumentException.class,
                () -> manager.generateRobocopyScripts(missing)
        );

        assertTrue(ex.getMessage().contains("Auto-generated folder does not exist"));
    }
}
