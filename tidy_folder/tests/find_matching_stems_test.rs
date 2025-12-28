#[cfg(test)]
mod tests {
    use std::collections::HashMap;
    use std::ffi::OsString;
    use std::fs::File;
    use tempfile::TempDir;
    use tidy_folder::find_matching_stems;
    use walkdir::{DirEntry, WalkDir};

    // Helper function to create test files and return walkdir::DirEntry objects
    fn create_test_files(dir: &TempDir, filenames: &[&str]) -> Vec<DirEntry> {
        let mut entries = Vec::new();

        for filename in filenames {
            let file_path = dir.path().join(filename);
            File::create(&file_path).unwrap();
        }

        // Use WalkDir to get entries
        for entry in WalkDir::new(dir.path()).min_depth(1).max_depth(1) {
            let entry = entry.unwrap();
            if entry.file_type().is_file() {
                entries.push(entry);
            }
        }

        entries
    }

    fn extract_stem_pairs(results: &[(DirEntry, DirEntry)]) -> Vec<(&str, &str)> {
        results.iter().map(|(file, other_file)| {
            let stem = file.path().file_stem().unwrap().to_str().unwrap();
            let other_stem = other_file.path().file_stem().unwrap().to_str().unwrap();
            (stem, other_stem)
        }).collect()
    }

    #[test]
    fn test_basic_prefix_match() {
        let temp_dir = TempDir::new().unwrap();
        let files = create_test_files(&temp_dir, &["data.txt", "data_backup.txt"]);

        let mut dirs_and_files = HashMap::new();
        dirs_and_files.insert(OsString::from("test_dir"), files);

        let results = find_matching_stems(dirs_and_files);

        assert_eq!(results.len(), 1);
        assert_eq!(results[0].0.file_name().to_str().unwrap(), "data.txt");
        assert_eq!(results[0].1.file_name().to_str().unwrap(), "data_backup.txt");
    }

    #[test]
    fn test_multiple_matches() {
        let temp_dir = TempDir::new().unwrap();
        let files = create_test_files(
            &temp_dir,
            &["report.csv", "report_2024.csv", "report_final.csv"],
        );

        let mut dirs_and_files = HashMap::new();
        dirs_and_files.insert(OsString::from("test_dir"), files);

        let results = find_matching_stems(dirs_and_files);

        assert_eq!(results.len(), 2);

        let found_pairs = extract_stem_pairs(&results);

        assert!(found_pairs.contains(&("report", "report_2024")));
        assert!(found_pairs.contains(&("report", "report_final")));
    }

    #[test]
    fn test_different_extensions_no_match() {
        let temp_dir = TempDir::new().unwrap();
        let files = create_test_files(&temp_dir, &["data.txt", "data_backup.csv"]);

        let mut dirs_and_files = HashMap::new();
        dirs_and_files.insert(OsString::from("test_dir"), files);

        let results = find_matching_stems(dirs_and_files);

        // Different extensions should not match
        assert_eq!(results.len(), 0);
    }

    #[test]
    fn test_no_prefix_relationship() {
        let temp_dir = TempDir::new().unwrap();
        let files = create_test_files(&temp_dir, &["apple.txt", "banana.txt", "cherry.txt"]);

        let mut dirs_and_files = HashMap::new();
        dirs_and_files.insert(OsString::from("test_dir"), files);

        let results = find_matching_stems(dirs_and_files);

        // No files have prefix relationships
        assert_eq!(results.len(), 0);
    }

    #[test]
    fn test_nested_prefixes() {
        let temp_dir = TempDir::new().unwrap();
        let files = create_test_files(
            &temp_dir,
            &["a.txt", "ab.txt", "abc.txt"]
        );

        let mut dirs_and_files = HashMap::new();
        dirs_and_files.insert(OsString::from("test_dir"), files);

        let results = find_matching_stems(dirs_and_files);

        assert_eq!(results.len(), 3);

        let found_pairs = extract_stem_pairs(&results);

        assert!(found_pairs.contains(&("a", "ab")));
        assert!(found_pairs.contains(&("a", "abc")));
        assert!(found_pairs.contains(&("ab", "abc")));
    }

    #[test]
    fn test_empty_input() {
        let dirs_and_files: HashMap<OsString, Vec<DirEntry>> = HashMap::new();

        let results = find_matching_stems(dirs_and_files);

        assert_eq!(results.len(), 0);
    }

    #[test]
    fn test_single_file() {
        let temp_dir = TempDir::new().unwrap();
        let files = create_test_files(&temp_dir, &["lonely.txt"]);

        let mut dirs_and_files = HashMap::new();
        dirs_and_files.insert(OsString::from("test_dir"), files);

        let results = find_matching_stems(dirs_and_files);

        assert_eq!(results.len(), 0);
    }

    #[test]
    fn test_files_without_extensions() {
        let temp_dir = TempDir::new().unwrap();
        let files = create_test_files(&temp_dir, &["README", "README_backup"]);

        let mut dirs_and_files = HashMap::new();
        dirs_and_files.insert(OsString::from("test_dir"), files);

        let results = find_matching_stems(dirs_and_files);

        // Files without extensions should be found
        assert_eq!(results.len(), 1);
    }

    #[test]
    fn test_multiple_directories() {
        let temp_dir1 = TempDir::new().unwrap();
        let temp_dir2 = TempDir::new().unwrap();

        let files1 = create_test_files(&temp_dir1, &["log.txt", "log_old.txt"]);
        let files2 = create_test_files(&temp_dir2, &["config.json", "config_backup.json"]);

        let mut dirs_and_files = HashMap::new();
        dirs_and_files.insert(OsString::from("dir1"), files1);
        dirs_and_files.insert(OsString::from("dir2"), files2);

        let results = find_matching_stems(dirs_and_files);

        // Should find one match in each directory
        assert_eq!(results.len(), 2);
    }

    #[test]
    fn test_underscore_and_dash_prefixes() {
        let temp_dir = TempDir::new().unwrap();
        let files = create_test_files(
            &temp_dir,
            &["project.rs", "project_test.rs", "project-backup.rs"]
        );

        let mut dirs_and_files = HashMap::new();
        dirs_and_files.insert(OsString::from("test_dir"), files);

        let results = find_matching_stems(dirs_and_files);

        assert_eq!(results.len(), 2);

        let found_pairs = extract_stem_pairs(&results);

        assert!(found_pairs.contains(&("project", "project_test")));
        assert!(found_pairs.contains(&("project", "project-backup")));
    }
}
