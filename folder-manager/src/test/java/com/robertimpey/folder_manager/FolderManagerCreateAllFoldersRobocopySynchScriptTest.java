package com.robertimpey.folder_manager;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertTrue;

public class FolderManagerCreateAllFoldersRobocopySynchScriptTest {
    @TempDir
    Path tempDir;

    @Test
    void createAllFoldersRobocopySynchScript_generatesCorrectContent() throws Exception {
        Path scriptPath = tempDir.resolve("all.ps1");
        Path sourcePath = Paths.get("C:\\Source");
        Path destinationPath = Paths.get("D:\\Destination");
        List<String> folders = List.of("Folder1", "Folder2");

        FolderManager.createAllFoldersRobocopySynchScript(folders, scriptPath, sourcePath, destinationPath);

        assertTrue(Files.exists(scriptPath));
        List<String> lines = Files.readAllLines(scriptPath);

        assertTrue(lines.stream().anyMatch(l -> l.contains("$folders = \"Folder1\", \"Folder2\"")));
        assertTrue(lines.stream().anyMatch(l -> l.contains("$src = \"" + sourcePath.toAbsolutePath() + "\"")));
        assertTrue(lines.stream().anyMatch(l -> l.contains("$dst = \"" + destinationPath.toAbsolutePath() + "\"")));
        assertTrue(lines.stream().anyMatch(l -> l.contains("foreach ($folder in $folders)")));
    }
}
