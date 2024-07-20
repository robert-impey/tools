use std::collections::HashMap;
use clap::{Parser};
use walkdir::{DirEntry, WalkDir};

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
struct Cli {
    directory: Option<String>,
}

fn main() {
    let cli = Cli::parse();

    if let Some(name) = cli.directory.as_deref() {
        let mut dirs_and_files: HashMap<String, Vec<DirEntry>> = HashMap::new();

        for entry in WalkDir::new(name).into_iter().filter_map(|e| e.ok()) {
            let path = entry.path();

            if path.is_file() {
                let parent = path.parent().unwrap();
                let parent_str = parent.to_str().unwrap().to_string();

                dirs_and_files.entry(parent_str)
                    .or_insert_with(Vec::new)
                    .push(entry);
            }
        }

        let mut matching_stems: Vec<(DirEntry, DirEntry)> = Vec::new();
        for (_dir, files) in dirs_and_files {
            for file in files.iter().cloned() {
                let path = file.path();

                match path.file_stem() {
                    Some(file_stem) =>
                        {
                            let file_stem_str = file_stem.to_str().unwrap();

                            match path.extension() {
                                Some(extension) => {
                                    let extension_string = extension.to_str().unwrap();

                                    for other_file in files.iter().cloned() {
                                        let other_path = other_file.path();

                                        match other_path.file_stem() {
                                            Some(other_file_stem) =>
                                                {
                                                    let other_file_stem_string = other_file_stem.to_str().unwrap();

                                                    match other_path.extension() {
                                                        Some(other_extension) => {
                                                            let other_extension_str = other_extension.to_str().unwrap();

                                                            if extension_string == other_extension_str {
                                                                if other_file_stem_string != file_stem_str {
                                                                    if other_file_stem_string.starts_with(file_stem_str) {
                                                                        matching_stems.push((file.clone(), other_file));
                                                                    }
                                                                }
                                                            }
                                                        }
                                                        None => ()
                                                    }
                                                }
                                            None => ()
                                        }
                                    }
                                }
                                None => ()
                            }
                        }
                    None => ()
                }
            }
        }

        if !matching_stems.is_empty() {
            println!("Value for directory: {name}");
            println!("Matching stems:");
            for (file, other_file) in matching_stems {
                println!("Matching stems in {}", file.path().parent().unwrap().display());
                println!("\t{}", file.file_name().to_str().unwrap());
                println!("\t{}", other_file.file_name().to_str().unwrap());
            }
        }
    }
}
