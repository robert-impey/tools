package com.robertimpey.folder_manager;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.io.IOException;
import java.io.PrintWriter;
import java.io.StringWriter;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
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
        assertEquals(List.of("C:/Data", "D:/Archive"), manager.locations());
        assertEquals(List.of("logs", "configs"), manager.folders());
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

        assertEquals(List.of("C:/RealPath"), manager.locations());
        assertEquals(List.of("folderA", "folderB"), manager.folders());
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

    @Test
    void getCleanLocationName_removesSpecialCharacters() {
        assertEquals("C_Data", FolderManager.getCleanLocationName("C:\\Data"));
        assertEquals("D_Archive", FolderManager.getCleanLocationName("D:/Archive/"));
        assertEquals("E_My_Documents", FolderManager.getCleanLocationName("E:\\My Documents"));
    }

    @Test
    void getScriptPath_buildsCorrectPath() {
        Path autoGen = Paths.get("C:\\AutoGen");
        Path expected = autoGen.resolve("C_Data").resolve("D_Backup");
        assertEquals(expected, FolderManager.getScriptPath(autoGen, "C:\\Data", "D:\\Backup"));
    }

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
        assertTrue(lines.stream().anyMatch(l -> l.contains("foreach ($folder in $folders)")));
    }

    @Test
    void listManagedFolders_listsExistingFolders() throws Exception {
        Path loc1 = tempDir.resolve("Loc1");
        Path loc2 = tempDir.resolve("Loc2");
        Files.createDirectories(loc1.resolve("FolderA"));
        Files.createDirectories(loc2.resolve("FolderB"));

        FolderManager manager = new FolderManager(
                List.of(loc1.toString(), loc2.toString()),
                List.of("FolderA", "FolderB", "FolderC")
        );

        StringWriter sw = new StringWriter();
        PrintWriter pw = new PrintWriter(sw);
        manager.listManagedFolders(pw);
        pw.flush();

        String output = sw.toString();
        assertTrue(output.contains("Loc1"));
        assertTrue(output.contains(loc1.resolve("FolderA").toAbsolutePath().toString()));
        assertTrue(output.contains("Loc2"));
        assertTrue(output.contains(loc2.resolve("FolderB").toAbsolutePath().toString()));
        assertFalse(output.contains("FolderC"));
    }

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