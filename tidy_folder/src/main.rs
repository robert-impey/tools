use std::{collections::HashMap, ffi::OsString};

use clap::Parser;
use walkdir::{DirEntry, WalkDir};
use itertools::Itertools; // Add `itertools = "0.13"` to Cargo.toml

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
struct Cli {
    directory: Option<String>,
}

fn main() {
    let cli = Cli::parse();

    if let Some(name) = cli.directory.as_deref() {
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

        let matching_stems: Vec<(DirEntry, DirEntry)> = dirs_and_files
            .into_iter()
            .flat_map(|(_dir, files)| {
                files
                    .into_iter()
                    .tuple_combinations() // Generates all unique pairs (a, b) where a != b
                    .filter_map(|(file, other_file)| {
                        let path = file.path();
                        let file_stem = path.file_stem()?.to_str()?;
                        let extension = path.extension()?.to_str()?;

                        let other_path = other_file.path();
                        let other_stem = other_path.file_stem()?.to_str()?;
                        let other_extension = other_path.extension()?.to_str()?;

                        if extension == other_extension
                            && other_stem.starts_with(file_stem)
                        {
                            Some((file, other_file))
                        } else {
                            None
                        }
                    })
                    .collect::<Vec<_>>()
            })
            .collect();

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
}
