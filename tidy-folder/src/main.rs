use std::collections::HashMap;
use clap::{Parser};
use walkdir::{DirEntry, WalkDir};

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
struct Cli {
    /// Optional type of file to search for
    file_type: Option<String>,

    directory: Option<String>,
}

fn main() {
    let cli = Cli::parse();

    if let Some(name) = cli.file_type.as_deref() {
        println!("Value for file_type: {name}");
    }

    if let Some(name) = cli.directory.as_deref() {
        println!("Value for directory: {name}");

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

        for (dir, files) in dirs_and_files {
            println!("Dir: {}", dir);
            for file in files.iter().cloned() {
                let path = file.path();
                println!("\tFile: {}", path.display());

                match path.file_stem() {
                    Some(file_stem) =>
                        {
                            println!("\t\tFile stem: {:?}", file_stem);

                            match path.extension() {
                                Some(extension) => {
                                    println!("\t\tExtension: {:?}", extension);

                                    for other_file in files.iter().cloned() {
                                        let other_path = other_file.path();
                                        println!("\tOther File: {}", other_path.display());

                                        match other_path.file_stem() {
                                            Some(other_file_stem) =>
                                                {
                                                    println!("\t\t\tOther File stem: {:?}", other_file_stem);

                                                    match other_path.extension() {
                                                        Some(other_extension) => {
                                                            println!("\t\t\tOther Extension: {:?}", other_extension);
                                                        },
                                                        None => println!("\t\tNo other extension")
                                                    }
                                                },
                                            None => println!("\t\tNo other file stem")
                                        }
                                    }
                                },
                                None => println!("\t\tNo extension")
                            }
                        },
                    None => println!("\t\tNo file stem")
                }
            }
        }
    }
}
