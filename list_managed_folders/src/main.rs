use chrono::{DateTime, Utc};
use clap::Parser;
use std::fs;
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
            println!("locations_path: {}", locations_path.display());
            println!("folders_path: {}", folders_path.display());

            let contents =
                make_managed_folders_file_contents(locations_path.into(), folders_path.into());

            if cli.write {
                if let Some(managed_folders_file_path) = cli.managed_folders_file.as_deref() {
                    println!("Writing to: {}", managed_folders_file_path.display());

                    let mut file_contents = String::new();

                    file_contents.push_str("# AUTOGEN'D!\n# ");

                    let now: DateTime<Utc> = Utc::now();

                    file_contents.push_str(now.to_rfc2822().as_str());

                    file_contents.push_str("\n");
                    file_contents.push_str("# DO NOT EDIT!\n");

                    file_contents.push_str(&contents);

                    let _ = fs::write(managed_folders_file_path, file_contents);
                }
            } else {
                println!("{}", contents);
            }
        }
    }
}

fn make_managed_folders_file_contents(
    locations_path: Box<Path>,
    folders_path: Box<Path>,
) -> String {
    let mut output = String::new();

    if let Some(locations_contents) = fs::read_to_string(locations_path).ok() {
        if let Some(folders_contents) = fs::read_to_string(folders_path).ok() {
            for location_line in locations_contents.lines() {
                let mut first = true;

                for folders_line in folders_contents.lines() {
                    let managed_folder_path = Path::new(location_line).join(folders_line);

                    if managed_folder_path.exists() {
                        if let Some(managed_folder_str) = managed_folder_path.to_str() {
                            if first {
                                output.push_str("\n");
                                first = false;
                            }
                            output.push_str(managed_folder_str);
                            output.push_str("\n");
                        }
                    }
                }
            }
        }
    }

    output
}
