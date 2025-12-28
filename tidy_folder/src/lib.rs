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

pub fn find_matching_stems(
    dirs_and_files: HashMap<OsString, Vec<DirEntry>>,
) -> Vec<(DirEntry, DirEntry)> {
    let mut matching_stems: Vec<(DirEntry, DirEntry)> = Vec::new();

    for (_dir, files) in dirs_and_files {
        for file in files.iter().cloned() {
            let path = file.path();

            let file_stem = match path.file_stem().and_then(|s| s.to_str()) {
                Some(s) => s,
                None => continue,
            };

            let extension = match path.extension().and_then(|s| s.to_str()) {
                Some(ext) => ext,
                None => continue,
            };

            for other_file in files.iter().cloned() {
                let other_path = other_file.path();

                let other_stem = match other_path.file_stem().and_then(|s| s.to_str()) {
                    Some(s) => s,
                    None => continue,
                };

                let other_extension = match other_path.extension().and_then(|s| s.to_str()) {
                    Some(ext) => ext,
                    None => continue,
                };

                if extension != other_extension {
                    continue;
                }

                if other_stem == file_stem {
                    continue;
                }

                if other_stem.starts_with(file_stem) {
                    matching_stems.push((file.clone(), other_file));
                }
            }
        }
    }

    matching_stems
}
