// From https://docs.rs/clap/latest/clap/_derive/_tutorial/chapter_0/index.html

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

        for entry in WalkDir::new(name).into_iter().filter_map(|e| e.ok()) {
            println!("{}", entry.path().display());
        }
    }
}
