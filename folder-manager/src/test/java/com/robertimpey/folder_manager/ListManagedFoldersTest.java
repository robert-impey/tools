package com.robertimpey.folder_manager;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.io.IOException;
import java.io.PrintWriter;
import java.io.StringWriter;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class ListManagedFoldersTest {
    @TempDir
    Path tempDir;

    @Test
    void listManagedFolders_shouldPrintExistingDirectoriesAndIgnoreMissing(@TempDir Path tempDir) throws IOException {
        // Given
        Path location1 = tempDir.resolve("loc1");
        Path location2 = tempDir.resolve("loc2");

        Files.createDirectories(location1.resolve("folder1")); // Valid
        Files.createDirectories(location1.resolve("folder2")); // Valid
        Files.createDirectories(location2.resolve("folder1")); // Valid (different location)

        // Create a file posing as a folder (should be ignored)
        Files.createDirectories(location2);
        Files.createFile(location2.resolve("folder2"));

        List<String> locations = Arrays.asList(
                location1.toAbsolutePath().toString(),
                location2.toAbsolutePath().toString()
        );
        List<String> folders = Arrays.asList("folder1", "folder2", "nonExistent");

        FolderManager folderManager = new FolderManager(locations, folders);
        StringWriter stringWriter = new StringWriter();
        PrintWriter printWriter = new PrintWriter(stringWriter);

        // When
        folderManager.listManagedFolders(printWriter);

        // Then
        String output = stringWriter.toString();
        String expectedPath1 = location1.resolve("folder1").toAbsolutePath().toString();
        String expectedPath2 = location1.resolve("folder2").toAbsolutePath().toString();
        String expectedPath3 = location2.resolve("folder1").toAbsolutePath().toString();

        assertTrue(output.contains(expectedPath1), "Output should contain loc1/folder1");
        assertTrue(output.contains(expectedPath2), "Output should contain loc1/folder2");
        assertTrue(output.contains(expectedPath3), "Output should contain loc2/folder1");

        assertFalse(output.contains("nonExistent"), "Output should not contain non-existent folders");
        assertFalse(output.contains(location2.resolve("folder2").toAbsolutePath().toString()), "Output should not contain files that match folder names");
    }

    @Test
    void listManagedFolders_shouldSeparateLocationsWithEmptyLine(@TempDir Path tempDir) throws IOException {
        // Given
        Path location1 = tempDir.resolve("loc1");
        Path location2 = tempDir.resolve("loc2");

        Files.createDirectories(location1.resolve("folderA"));
        Files.createDirectories(location2.resolve("folderA"));

        List<String> locations = Arrays.asList(
                location1.toAbsolutePath().toString(),
                location2.toAbsolutePath().toString()
        );
        List<String> folders = Collections.singletonList("folderA");

        FolderManager folderManager = new FolderManager(locations, folders);
        StringWriter stringWriter = new StringWriter();
        PrintWriter printWriter = new PrintWriter(stringWriter);

        // When
        folderManager.listManagedFolders(printWriter);

        // Then
        String output = stringWriter.toString().trim();
        String[] lines = output.replace("\r\n", "\n").split("\n");

        // Expecting:
        // Path1
        // <empty line>
        // Path2
        assertEquals(3, lines.length, "Should have 3 lines (path, empty, path)");
        assertEquals(location1.resolve("folderA").toAbsolutePath().toString(), lines[0]);
        assertTrue(lines[1].isEmpty(), "Second line should be empty separator");
        assertEquals(location2.resolve("folderA").toAbsolutePath().toString(), lines[2]);
    }

    @Test
    void listManagedFolders_shouldIgnoreSymbolicLinks(@TempDir Path tempDir) throws IOException {
        // Given
        Path source = tempDir.resolve("real_source");
        Path location = tempDir.resolve("location");
        Files.createDirectories(source.resolve("myFolder"));
        Files.createDirectories(location);

        // Create a symlink in 'location' pointing to 'real_source/myFolder'
        Path symlink = location.resolve("myFolder");
        try {
            Files.createSymbolicLink(symlink, source.resolve("myFolder"));
        } catch (UnsupportedOperationException e) {
            // Skip test on OS that doesn't support symlinks or requires admin rights (e.g. strict Windows env)
            System.out.println("Skipping symlink test: operation not supported.");
            return;
        } catch (IOException e) {
            // Windows often throws AccessDeniedException for symlinks if not admin
            System.out.println("Skipping symlink test: " + e.getMessage());
            return;
        }

        List<String> locations = Collections.singletonList(location.toAbsolutePath().toString());
        List<String> folders = Collections.singletonList("myFolder");

        FolderManager folderManager = new FolderManager(locations, folders);
        StringWriter stringWriter = new StringWriter();
        PrintWriter printWriter = new PrintWriter(stringWriter);

        // When
        folderManager.listManagedFolders(printWriter);

        // Then
        String output = stringWriter.toString();
        assertTrue(output.isEmpty(), "Output should be empty because symlinks must be ignored");
    }
}
