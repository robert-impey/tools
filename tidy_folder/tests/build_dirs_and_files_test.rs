#[cfg(test)]
mod tests {
    use std::fs;
    use tempfile::TempDir;
    use tidy_folder::build_dirs_and_files;

    #[test]
    fn test_empty_directory() {
        let temp_dir = TempDir::new().unwrap();
        let result = build_dirs_and_files(temp_dir.path().to_str().unwrap());
        assert!(
            result.is_empty(),
            "Empty directory should produce empty HashMap"
        );
    }

    #[test]
    fn test_single_file_in_subdirectory() {
        let temp_dir = TempDir::new().unwrap();
        let sub_dir = temp_dir.path().join("subdir");
        fs::create_dir(&sub_dir).unwrap();
        fs::write(sub_dir.join("file.txt"), "content").unwrap();

        let result = build_dirs_and_files(temp_dir.path().to_str().unwrap());

        assert_eq!(result.len(), 1, "Should have one directory entry");
        let files = result.get(sub_dir.as_os_str()).unwrap();
        assert_eq!(files.len(), 1, "Should have one file in subdir");
    }

    #[test]
    fn test_multiple_files_same_directory() {
        let temp_dir = TempDir::new().unwrap();
        let sub_dir = temp_dir.path().join("subdir");
        fs::create_dir(&sub_dir).unwrap();
        fs::write(sub_dir.join("file1.txt"), "content1").unwrap();
        fs::write(sub_dir.join("file2.txt"), "content2").unwrap();
        fs::write(sub_dir.join("file3.txt"), "content3").unwrap();

        let result = build_dirs_and_files(temp_dir.path().to_str().unwrap());

        let files = result.get(sub_dir.as_os_str()).unwrap();
        assert_eq!(files.len(), 3, "Should have three files in subdir");
    }

    #[test]
    fn test_nested_directories() {
        let temp_dir = TempDir::new().unwrap();
        let level1 = temp_dir.path().join("level1");
        let level2 = level1.join("level2");
        fs::create_dir_all(&level2).unwrap();

        fs::write(level1.join("file1.txt"), "content1").unwrap();
        fs::write(level2.join("file2.txt"), "content2").unwrap();

        let result = build_dirs_and_files(temp_dir.path().to_str().unwrap());

        assert_eq!(result.len(), 2, "Should have two directory entries");
        assert_eq!(result.get(level1.as_os_str()).unwrap().len(), 1);
        assert_eq!(result.get(level2.as_os_str()).unwrap().len(), 1);
    }

    #[test]
    fn test_multiple_directories_with_files() {
        let temp_dir = TempDir::new().unwrap();
        let src_dir = temp_dir.path().join("src");
        let tests_dir = temp_dir.path().join("tests");
        fs::create_dir(&src_dir).unwrap();
        fs::create_dir(&tests_dir).unwrap();

        fs::write(src_dir.join("main.rs"), "fn main() {}").unwrap();
        fs::write(src_dir.join("lib.rs"), "pub fn lib() {}").unwrap();
        fs::write(tests_dir.join("test1.rs"), "#[test]").unwrap();

        let result = build_dirs_and_files(temp_dir.path().to_str().unwrap());

        assert_eq!(result.len(), 2, "Should have two directory entries");
        assert_eq!(result.get(src_dir.as_os_str()).unwrap().len(), 2);
        assert_eq!(result.get(tests_dir.as_os_str()).unwrap().len(), 1);
    }

    #[test]
    fn test_file_in_root_directory() {
        let temp_dir = TempDir::new().unwrap();
        fs::write(temp_dir.path().join("root_file.txt"), "content").unwrap();

        let result = build_dirs_and_files(temp_dir.path().to_str().unwrap());

        // The root file should be included with its parent directory as the key
        assert_eq!(result.len(), 1, "Should have one directory entry for root");
    }

    #[test]
    fn test_only_directories_no_files() {
        let temp_dir = TempDir::new().unwrap();
        let sub1 = temp_dir.path().join("sub1");
        let sub2 = temp_dir.path().join("sub2");
        fs::create_dir(&sub1).unwrap();
        fs::create_dir(&sub2).unwrap();

        let result = build_dirs_and_files(temp_dir.path().to_str().unwrap());

        assert!(
            result.is_empty(),
            "Directories without files should produce empty HashMap"
        );
    }

    #[test]
    fn test_mixed_file_types() {
        let temp_dir = TempDir::new().unwrap();
        let dir = temp_dir.path().join("mixed");
        fs::create_dir(&dir).unwrap();

        fs::write(dir.join("file.txt"), "text").unwrap();
        fs::write(dir.join("data.json"), "{}").unwrap();
        fs::write(dir.join("script.sh"), "#!/bin/bash").unwrap();

        let result = build_dirs_and_files(temp_dir.path().to_str().unwrap());

        let files = result.get(dir.as_os_str()).unwrap();
        assert_eq!(files.len(), 3, "Should include all file types");
    }

    #[test]
    fn test_nonexistent_directory() {
        let result = build_dirs_and_files("/nonexistent/path/that/does/not/exist");
        assert!(
            result.is_empty(),
            "Nonexistent path should produce empty HashMap"
        );
    }
}
