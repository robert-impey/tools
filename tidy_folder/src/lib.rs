use chrono::Local;
use std::collections::HashMap;
use std::ffi::OsString;
use std::fs::File;
use std::io;
use std::io::Write;
use std::path::{Path, PathBuf};
use walkdir::{DirEntry, WalkDir};

pub fn process_directory(dir: &str, logs_dir: Option<&Path>) -> anyhow::Result<()> {
    // 1. Determine our sinks (Log vs Stdout)
    let mut log_sink: Box<dyn Write> = match logs_dir {
        Some(path) => {
            let safe = dir.replace(['/', '\\', ':'], "_");
            let timestamp = get_log_time();
            let log_path = path.join(format!("{timestamp}-search-{safe}.log"));
            Box::new(File::create(log_path)?)
        }
        None => Box::new(io::stdout()), // Fallback to stdout
    };

    // 2. Determine our error sink (File vs Stderr)
    let mut err_sink: Box<dyn Write> = match logs_dir {
        Some(path) => {
            let safe = dir.replace(['/', '\\', ':'], "_");
            let timestamp = get_log_time();
            let err_path = path.join(format!("{timestamp}-search-{safe}.err"));
            Box::new(File::create(err_path)?)
        }
        None => Box::new(io::stderr()), // Fallback to stderr
    };

    // 3. Wrap logic to use these sinks
    match (|| {
        let dirs_and_files = build_dirs_and_files(dir);
        let matching_stems = find_matching_stems(dirs_and_files);

        // Ensure print_matching_stems accepts &mut dyn Write
        print_matching_stems(&mut *log_sink, dir, &matching_stems)?;
        Ok::<_, anyhow::Error>(())
    })() {
        Ok(_) => {
            // We use eprintln so this status message doesn't get
            // mixed into the data if log_sink is currently stdout
            if logs_dir.is_some() {
                println!("OK: processed {dir}");
            }
        }
        Err(e) => {
            writeln!(err_sink, "ERROR processing {dir}: {e}")?;
        }
    }

    Ok(())
}

fn get_log_time() -> String {
    Local::now().format("%Y-%m-%d_%H.%M.%S").to_string()
}

fn print_matching_stems<W: Write>(
    mut out: W,
    name: &str,
    matching_stems: &[(DirEntry, DirEntry)],
) -> io::Result<()> {
    if matching_stems.is_empty() {
        return Ok(());
    }

    // A simple visual separator makes stdout much more readable
    writeln!(out, "\n--- Results for: {} ---", name)?;

    for (file, other_file) in matching_stems {
        // parent().unwrap() is risky if the path is the root;
        // display() handles the rest nicely.
        let parent = file
            .path()
            .parent()
            .map(|p| p.display().to_string())
            .unwrap_or_default();

        writeln!(out, "{}", parent)?;
        writeln!(out, "\t{}", file.file_name().to_string_lossy())?;
        writeln!(out, "\t{}", other_file.file_name().to_string_lossy())?;
    }

    Ok(())
}

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

            // Treat missing extension as empty string
            let extension = path.extension().and_then(|s| s.to_str()).unwrap_or("");

            for other_file in files.iter().cloned() {
                let other_path = other_file.path();

                let other_stem = match other_path.file_stem().and_then(|s| s.to_str()) {
                    Some(s) => s,
                    None => continue,
                };

                // Treat missing extension as empty string
                let other_extension = other_path
                    .extension()
                    .and_then(|s| s.to_str())
                    .unwrap_or("");

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
