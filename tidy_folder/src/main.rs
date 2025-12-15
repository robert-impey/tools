use std::{collections::HashMap, ffi::OsString};

use clap::Parser;
use walkdir::{DirEntry, WalkDir};

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
struct Cli {
    directory: Option<String>,
}

fn main() {
    let cli = Cli::parse();

    if let Some(name) = cli.directory.as_deref() {
        let dirs_and_files =build_dirs_and_files(name);

        let matching_stems = find_matching_stems(dirs_and_files);

        print_matching_stems(name, matching_stems);
    }
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
    dirs_and_files: HashMap<OsString, Vec<DirEntry>>
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

fn print_matching_stems(name: &str, matching_stems: Vec<(DirEntry, DirEntry)>) {
    if !matching_stems.is_empty() {
        println!("Value for directory: {name}");
        println!("Matching stems:");
        for (file, other_file) in matching_stems {
            println!(
                "Matching stems in {}",
                file.path().parent().unwrap().display()
            );
            println!("\t{}", file.file_name().to_str().unwrap());
            println!("\t{}", other_file.file_name().to_str().unwrap());
        }
    }
}
