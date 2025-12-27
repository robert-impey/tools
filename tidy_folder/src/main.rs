use std::{collections::HashMap, ffi::OsString};

use chrono::Local;
use clap::{Parser, Subcommand};
use std::io::{self, Write};
use std::path::Path;
use std::path::PathBuf;
use walkdir::{DirEntry, WalkDir};
use tidy_folder::read_directories;

#[derive(Parser)]
#[command(author, version, about)]
pub struct Cli {
    #[command(subcommand)]
    pub command: Commands,
}

#[derive(Subcommand)]
pub enum Commands {
    /// Search a single directory
    Search {
        /// Directory to search
        directory: PathBuf,

        /// Directory where logs should be written
        #[arg(long = "logs-dir")]
        logs_dir: PathBuf,
    },

    /// Search multiple directories listed in a text file
    SearchFrom {
        /// File containing directories to process
        directories_file: PathBuf,

        /// Directory where logs should be written
        #[arg(long = "logs-dir")]
        logs_dir: PathBuf,
    },
}

fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();

    match cli.command {
        Commands::Search {
            directory,
            logs_dir,
        } => {
            process_directory(&directory.to_string_lossy(), &logs_dir)?;
        }

        Commands::SearchFrom {
            directories_file,
            logs_dir,
        } => {
            let dirs = read_directories(&directories_file)?;
            for dir in dirs {
                process_directory(&dir, &logs_dir)?;
            }
        }
    }

    Ok(())
}

fn process_directory(dir: &str, logs_dir: &Path) -> anyhow::Result<()> {
    // Build log-safe filename
    let safe = dir.replace('/', "_").replace('\\', "_").replace(':', "");

    let timestamp = get_log_time(); // your existing function
    let log_path = logs_dir.join(format!("{timestamp}-search-{safe}.log"));
    let err_path = logs_dir.join(format!("{timestamp}-search-{safe}.err"));

    let mut log_file = std::fs::File::create(log_path)?;
    let mut err_file = std::fs::File::create(err_path)?;

    // Capture stdout/stderr manually
    use std::io::Write;

    // Wrap your logic so you can write logs deterministically
    match (|| {
        let dirs_and_files = build_dirs_and_files(dir);
        let matching_stems = find_matching_stems(dirs_and_files);
        print_matching_stems(&mut log_file, dir, &matching_stems)?;
        Ok::<_, anyhow::Error>(())
    })() {
        Ok(_) => {
            println!("OK: processed {dir}");
        }
        Err(e) => {
            writeln!(err_file, "ERROR processing {dir}: {e}")?;
        }
    }

    Ok(())
}

fn get_log_time() -> String {
    Local::now().format("%Y-%m-%d_%H.%M.%S").to_string()
}

fn build_dirs_and_files(name: &str) -> HashMap<OsString, Vec<DirEntry>> {
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

fn find_matching_stems(
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

fn print_matching_stems<W: Write>(
    mut out: W,
    name: &str,
    matching_stems: &[(DirEntry, DirEntry)],
) -> io::Result<()> {
    if matching_stems.is_empty() {
        return Ok(());
    }

    writeln!(out, "Value for directory: {name}")?;
    writeln!(out, "Matching stems:")?;

    for (file, other_file) in matching_stems {
        writeln!(
            out,
            "Matching stems in {}",
            file.path().parent().unwrap().display()
        )?;
        writeln!(out, "\t{}", file.file_name().to_string_lossy())?;
        writeln!(out, "\t{}", other_file.file_name().to_string_lossy())?;
    }

    Ok(())
}
