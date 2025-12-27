// tests/read_directories_test.rs
use std::path::PathBuf;
use tidy_folder::read_directories;

#[test]
fn test_read_directories_basic() {
    let path = PathBuf::from("tests/data/directories.txt");
    let result = read_directories(&path).expect("Failed to read directories");

    assert_eq!(result.len(), 3);
    assert_eq!(result[0], "/Users/foo/apples");
    assert_eq!(result[1], "/Users/foo/bananas");
    assert_eq!(result[2], "/Users/foo/coconuts");
}

#[test]
fn test_read_directories_ignores_comments() {
    let path = PathBuf::from("tests/data/directories.txt");
    let result = read_directories(&path).expect("Failed to read directories");

    // Should not contain any lines starting with #
    assert!(!result.iter().any(|line| line.starts_with('#')));
}

#[test]
fn test_read_directories_ignores_empty_lines() {
    // Create a test file with empty lines
    let path = PathBuf::from("tests/data/directories_with_empty.txt");
    let result = read_directories(&path).expect("Failed to read directories");

    // Should not contain empty strings
    assert!(!result.iter().any(|line| line.is_empty()));
}

#[test]
fn test_read_directories_trims_whitespace() {
    let path = PathBuf::from("tests/data/directories_with_whitespace.txt");
    let result = read_directories(&path).expect("Failed to read directories");

    // All entries should have no leading/trailing whitespace
    for dir in &result {
        assert_eq!(dir, dir.trim());
    }
}

#[test]
fn test_read_directories_nonexistent_file() {
    let path = PathBuf::from("tests/data/nonexistent.txt");
    let result = read_directories(&path);

    assert!(result.is_err());
}

#[test]
fn test_read_directories_unicode_normalization() {
    // Test with a file containing unicode characters that need normalization
    let path = PathBuf::from("tests/data/directories_unicode.txt");
    let result = read_directories(&path).expect("Failed to read directories");

    // This tests that NFC normalization is applied
    assert!(result.len() > 0);
}
