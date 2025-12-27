use std::{collections::HashMap, ffi::OsString};

use chrono::Local;
use clap::{Parser, Subcommand};
use std::io::{self, Write};
use std::path::Path;
use std::path::PathBuf;
use tidy_folder::{build_dirs_and_files, read_directories};
use walkdir::DirEntry;

use std::fs::File;

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

        /// Directory where logs should be written (Optional)
        #[arg(long = "logs-dir")]
        logs_dir: Option<PathBuf>, // Becomes --logs-dir <PATH>
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
            // logs_dir is Option<PathBuf>
            // .as_deref() converts Option<PathBuf> to Option<&Path>
            process_directory(&directory.to_string_lossy(), logs_dir.as_deref())?;
        }

        Commands::SearchFrom {
            directories_file,
            logs_dir,
        } => {
            // logs_dir is PathBuf
            let dirs = read_directories(&directories_file)?;
            for dir in dirs {
                // We wrap it in Some() to match the Option<&Path> signature
                process_directory(&dir, Some(&logs_dir))?;
            }
        }
    }

    Ok(())
}

fn process_directory(dir: &str, logs_dir: Option<&Path>) -> anyhow::Result<()> {
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
