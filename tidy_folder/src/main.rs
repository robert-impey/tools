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

        let mut matching_stems: Vec<(DirEntry, DirEntry)> = Vec::new();
        for (_dir, files) in dirs_and_files {
            for file in files.iter().cloned() {
                let path = file.path();

                if let Some(file_stem) = path.file_stem() {
                    match path.extension() {
                        Some(extension) => {
                            let extension_string = extension.to_str().unwrap();

                            for other_file in files.iter().cloned() {
                                let other_path = other_file.path();

                                let other_file_stem = other_path.file_stem().unwrap();

                                match other_path.extension() {
                                    Some(other_extension) => {
                                        let other_extension_str = other_extension.to_str().unwrap();

                                        if extension_string == other_extension_str {
                                            if other_file_stem != file_stem {
                                                if let Some(other_file_stem_str) =
                                                    other_file_stem.to_str()
                                                {
                                                    if let Some(file_stem_str) = file_stem.to_str()
                                                    {
                                                        if other_file_stem_str
                                                            .starts_with(file_stem_str)
                                                        {
                                                            matching_stems
                                                                .push((file.clone(), other_file));
                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    }
                                    None => (),
                                }
                            }
                        }
                        None => (),
                    }
                }
            }
        }

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
