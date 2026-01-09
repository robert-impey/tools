use clap::{Parser, Subcommand};
use std::path::PathBuf;
use tidy_folder::{search_directory, read_directories};

#[derive(Parser)]
#[command(author, version, about)]
pub struct Cli {
    #[command(subcommand)]
    pub command: Commands,
}

#[derive(Subcommand)]
pub enum Commands {
    /// Search a single directory
    Search {
        /// Directory to search
        directory: PathBuf,

        /// Directory where logs should be written (Optional)
        #[arg(long = "logs-dir")]
        logs_dir: Option<PathBuf>, // Becomes --logs-dir <PATH>
    },

    /// Search multiple directories listed in a text file
    SearchFrom {
        /// File containing directories to process
        directories_file: PathBuf,

        /// Directory where logs should be written
        #[arg(long = "logs-dir")]
        logs_dir: PathBuf,
    },
}

fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();

    match cli.command {
        Commands::Search {
            directory,
            logs_dir,
        } => {
            // logs_dir is Option<PathBuf>
            // .as_deref() converts Option<PathBuf> to Option<&Path>
            search_directory(&directory.to_string_lossy(), logs_dir.as_deref())?;
        }

        Commands::SearchFrom {
            directories_file,
            logs_dir,
        } => {
            // logs_dir is PathBuf
            let dirs = read_directories(&directories_file)?;
            for dir in dirs {
                // We wrap it in Some() to match the Option<&Path> signature
                search_directory(&dir, Some(&logs_dir))?;
            }
        }
    }

    Ok(())
}
