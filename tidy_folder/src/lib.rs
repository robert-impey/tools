use std::collections::HashMap;
use std::ffi::OsString;
use std::path::PathBuf;
use walkdir::{DirEntry, WalkDir};

pub fn read_directories(path: &PathBuf) -> std::io::Result<Vec<String>> {
    use std::io::{BufRead, BufReader};
    use unicode_normalization::UnicodeNormalization;

    let file = std::fs::File::open(path)?;
    let reader = BufReader::new(file);

    let mut dirs = Vec::new();

    for line in reader.lines() {
        let mut line = line?.trim().to_string();

        if line.is_empty() || line.starts_with('#') {
            continue;
        }

        line = line.nfc().collect();
        dirs.push(line);
    }

    Ok(dirs)
}

pub fn build_dirs_and_files(name: &str) -> HashMap<OsString, Vec<DirEntry>> {
    let mut dirs_and_files: HashMap<OsString, Vec<DirEntry>> = HashMap::new();

    for entry in WalkDir::new(name).into_iter().filter_map(|e| e.ok()) {
        if entry.file_type().is_dir() {
            continue;
        }

        let path = entry.path();

        if let Some(parent) = path.parent() {
            dirs_and_files
                .entry(parent.to_path_buf().into_os_string())
                .or_insert_with(Vec::new)
                .push(entry);
        }
    }

    dirs_and_files
}
