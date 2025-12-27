use std::path::PathBuf;

pub fn read_directories(path: &PathBuf) -> std::io::Result<Vec<String>> {
    use std::io::{BufRead, BufReader};
    use unicode_normalization::UnicodeNormalization;

    let file = std::fs::File::open(path)?;
    let reader = BufReader::new(file);

    let mut dirs = Vec::new();

    for line in reader.lines() {
        let mut line = line?.trim().to_string();

        if line.is_empty() || line.starts_with('#') {
            continue;
        }

        line = line.nfc().collect();
        dirs.push(line);
    }

    Ok(dirs)
}