package com.robertimpey.folder_manager;

import org.junit.jupiter.api.Test;

import java.nio.file.Path;
import java.nio.file.Paths;

import static org.junit.jupiter.api.Assertions.assertEquals;

class FolderManagerTest {
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
}