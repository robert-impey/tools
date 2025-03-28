use std::path::PathBuf;
use clap::{Parser};

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

    if let Some(folders_path) = cli.folders.as_deref()
    {
        println!("Value for folders_path: {}", folders_path.display());

        if let Some(locations_path) = cli.locations.as_deref()
        {
            println!("Value for locations_path: {}", locations_path.display());
            
            if cli.write {
                if let Some(managed_folders_file_path) = cli.managed_folders_file.as_deref()
                {
                    println!("Value for managed_folders_file_path: {}", managed_folders_file_path.display());
                }
            }
        }
    }
}
