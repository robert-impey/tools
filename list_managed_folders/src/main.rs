use clap::Parser;
use std::path::{Path, PathBuf};

#[derive(Parser)]
#[command(version, about, long_about = None)]
struct Cli {
    #[arg(short, long, value_name = "FOLDERS_FILE")]
    folders: Option<PathBuf>,

    #[arg(short, long, value_name = "LOCATIONS_FILE")]
    locations: Option<PathBuf>,

    #[arg(short, long, value_name = "MANAGED_FOLDERS_FILE")]
    managed_folders_file: Option<PathBuf>,

    #[arg(short, long)]
    write: bool,
}

fn main() {
    let cli = Cli::parse();

    if let Some(folders_path) = cli.folders.as_deref() {
        if let Some(locations_path) = cli.locations.as_deref() {
            if cli.write {
                if let Some(managed_folders_file_path) = cli.managed_folders_file.as_deref() {
                    write_managed_folders_file(
                        locations_path.into(),
                        folders_path.into(),
                        managed_folders_file_path.into(),
                    );
                }
            } else {
                print_managed_folders(locations_path.into(), folders_path.into());
            }
        }
    }
}

fn print_managed_folders(locations_path: Box<Path>, folders_path: Box<Path>) {
    println!("folders_path: {}", folders_path.display());
    println!("locations_path: {}", locations_path.display());
}

fn write_managed_folders_file(
    locations_path: Box<Path>,
    folders_path: Box<Path>,
    managed_folders_file_path: Box<Path>,
) {
    println!("folders_path: {}", folders_path.display());
    println!("locations_path: {}", locations_path.display());
    println!("Writing to: {}", managed_folders_file_path.display());
}
