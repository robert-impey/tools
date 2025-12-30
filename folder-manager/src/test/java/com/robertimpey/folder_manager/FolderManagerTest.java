package com.robertimpey.folder_manager;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class FolderManagerTest {

    @TempDir
    Path tempDir;

    private Path writeFile(String name, List<String> lines) throws IOException {
        Path file = tempDir.resolve(name);
        Files.write(file, lines);
        return file;
    }

    @Test
    void create_successfullyLoadsLocationsAndFolders() throws Exception {
        Path locationsFile = writeFile("locations.txt", List.of(
                "C:/Data",
                "D:/Archive"
        ));

        Path foldersFile = writeFile("folders.txt", List.of(
                "logs",
                "configs"
        ));

        FolderManager manager = FolderManager.create(locationsFile, foldersFile);

        assertNotNull(manager);
        assertEquals(List.of("C:/Data", "D:/Archive"), manager.getLocations());
        assertEquals(List.of("logs", "configs"), manager.getFolders());
    }

    @Test
    void create_ignoresBlankAndCommentLines() throws Exception {
        Path locationsFile = writeFile("locations.txt", List.of(
                "",
                "   ",
                "# comment",
                "C:/RealPath"
        ));

        Path foldersFile = writeFile("folders.txt", List.of(
                "# header",
                "folderA",
                "",
                "folderB"
        ));

        FolderManager manager = FolderManager.create(locationsFile, foldersFile);

        assertEquals(List.of("C:/RealPath"), manager.getLocations());
        assertEquals(List.of("folderA", "folderB"), manager.getFolders());
    }

    @Test
    void create_throwsIfLocationsEmpty() throws Exception {
        Path locationsFile = writeFile("locations.txt", List.of(
                "   ",
                "# comment"
        ));

        Path foldersFile = writeFile("folders.txt", List.of("folderA"));

        IllegalArgumentException ex = assertThrows(
                IllegalArgumentException.class,
                () -> FolderManager.create(locationsFile, foldersFile)
        );

        assertEquals("Locations cannot be null or empty", ex.getMessage());
    }

    @Test
    void create_throwsIfFoldersEmpty() throws Exception {
        Path locationsFile = writeFile("locations.txt", List.of("C:/Data"));

        Path foldersFile = writeFile("folders.txt", List.of(
                "",
                "# comment"
        ));

        IllegalArgumentException ex = assertThrows(
                IllegalArgumentException.class,
                () -> FolderManager.create(locationsFile, foldersFile)
        );

        assertEquals("Folders cannot be null or empty", ex.getMessage());
    }

    @Test
    void create_throwsIfFileDoesNotExist() {
        Path missing = tempDir.resolve("missing.txt");
        Path foldersFile = tempDir.resolve("folders.txt");

        assertDoesNotThrow(() -> Files.write(foldersFile, List.of("folder")));

        Exception ex = assertThrows(
                Exception.class,
                () -> FolderManager.create(missing, foldersFile)
        );

        assertTrue(ex.getMessage().contains("File does not exist"));
    }
}