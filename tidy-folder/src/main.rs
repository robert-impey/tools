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

        let mut matching_stems: Vec<(DirEntry, DirEntry)> = Vec::new();
        for (_dir, files) in dirs_and_files {
            //println!("Dir: {}", dir);
            for file in files.iter().cloned() {
                let path = file.path();
                //println!("\tFile: {}", path.display());

                match path.file_stem() {
                    Some(file_stem) =>
                        {
                            let file_stem_str = file_stem.to_str().unwrap();

                            //println!("\t\tFile stem: {}", file_stem_str);

                            match path.extension() {
                                Some(extension) => {
                                    let extension_string = extension.to_str().unwrap();

                                    //println!("\t\tExtension: {}", extension_string);

                                    for other_file in files.iter().cloned() {
                                        let other_path = other_file.path();
                                        //println!("\tOther File: {}", other_path.display());

                                        match other_path.file_stem() {
                                            Some(other_file_stem) =>
                                                {
                                                    let other_file_stem_string = other_file_stem.to_str().unwrap();

                                                    //println!("\t\t\tOther File stem: {}", other_file_stem_string);

                                                    match other_path.extension() {
                                                        Some(other_extension) => {
                                                            let other_extension_str = other_extension.to_str().unwrap();

                                                            //println!("\t\t\tOther Extension: {}", other_extension_str);

                                                            if extension_string == other_extension_str {
                                                                //println!("\t\t\t\tMatching extension");

                                                                if other_file_stem_string == file_stem_str {
                                                                    //println!("\t\t\t\t{} is {}. Skipping...", other_file_stem_string, file_stem_str);
                                                                } else {
                                                                    if other_file_stem_string.starts_with(file_stem_str) {
                                                                        //println!("\t\t\t\tFile stem match!");
                                                                        matching_stems.push((file.clone(), other_file));
                                                                    } else {
                                                                        //println!("\t\t\t\t{} does not start with {}", other_file_stem_string, file_stem_str);
                                                                    }
                                                                }
                                                            } else {
                                                                //println!("\t\t\t\t{} does not equal {}", extension_string, other_extension_str);
                                                            }
                                                        }
                                                        None => ()//println!("\t\tNo other extension")
                                                    }
                                                }
                                            None => () // println!("\t\tNo other file stem")
                                        }
                                    }
                                }
                                None => () //println!("\t\tNo extension")
                            }
                        }
                    None => () //println!("\t\tNo file stem")
                }
            }
        }

        println!("Matching stems:");
        for (file, other_file) in matching_stems {
            println!("Matching stems {} and {}", file.path().display(), other_file.path().display());
        }
    }
}
