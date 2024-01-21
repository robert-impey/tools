use std::collections::HashMap;
use clap::{Parser};
use walkdir::WalkDir;

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

        let mut dirs_and_files: HashMap<String, Vec<String>>  = HashMap::new();

        for entry in WalkDir::new(name).into_iter().filter_map(|e| e.ok()) {
            let path = entry.path();
            let path_str = path.to_str().unwrap().to_string();

            if path.is_file() {
                let parent = path.parent().unwrap();
                let parent_str = parent.to_str().unwrap().to_string();

                dirs_and_files.entry(parent_str)
                    .or_insert_with(Vec::new)
                    .push(path_str);
            }
        }

        for (dir, files) in dirs_and_files {
            println!("Dir: {}", dir);
            for file in files {
                println!("\tFile: {}", file);
            }
        }
    }
}
