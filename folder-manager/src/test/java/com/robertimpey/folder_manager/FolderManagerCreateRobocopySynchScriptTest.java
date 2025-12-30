package com.robertimpey.folder_manager;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertTrue;

public class FolderManagerCreateRobocopySynchScriptTest {
    @TempDir
    Path tempDir;

    @Test
    void createRobocopySynchScript_generatesCorrectContent() throws Exception {
        Path scriptPath = tempDir.resolve("script.ps1");
        Path sourcePath = Paths.get("C:\\Source");
        Path destinationPath = Paths.get("D:\\Destination");

        FolderManager.createRobocopySynchScript("MyFolder", scriptPath, sourcePath, destinationPath);

        assertTrue(Files.exists(scriptPath));
        List<String> lines = Files.readAllLines(scriptPath);

        assertTrue(lines.stream().anyMatch(l -> l.contains("$folder = \"MyFolder\"")));
        assertTrue(lines.stream().anyMatch(l -> l.contains("$src = \"C:\\Source\"")));
        assertTrue(lines.stream().anyMatch(l -> l.contains("$dst = \"D:\\Destination\"")));
        assertTrue(lines.stream().anyMatch(l -> l.contains("Synch $folder $src $dst $logged")));
    }
}
