## ffheic

A tiny Bash utility that batch‑converts **HEIC** images to **PNG** or **JPG** using **ffmpeg**.  
It works in any POSIX‑compatible shell (Linux/macOS, Git Bash on Windows, etc.).

---  

### Table of Contents
1. [Prerequisites](#prerequisites)  
2. [Installation](#installation)  
3. [Usage](#usage)  
4. [Options & Arguments](#options--arguments)  
5. [Examples](#examples)  
6. [How It Works](#how-it-works)  
7. [Troubleshooting](#troubleshooting)  

---  

## Prerequisites
| Requirement | Why it’s needed |
|-------------|-----------------|
| **Bash** (or any POSIX‑compatible shell) | Executes the script |
| **ffmpeg** (≥ 4.0) | Performs the actual image conversion |
| **Git Bash** (Windows only) | Provides a Bash environment on Windows |

Make sure `ffmpeg` is reachable from your `PATH`:

```bash
ffmpeg -version   # should print version information
```

If it isn’t installed, see the [ffmpeg download page](https://ffmpeg.org/download.html).

---  

## Installation

1. **Clone the repository (or copy the script)**  

   ```bash
   git clone https://github.com/yourname/ffheic.git
   cd ffheic
   ```

2. **Make the script executable**  

   ```bash
   chmod +x ffheic.sh
   ```

3 (optional). **Add the script to your PATH** for easy access:

```bash
# Example for a single‑user setup
mkdir -p "$HOME/.local/bin"
cp ffheic.sh "$HOME/.local/bin/ffheic"
# Ensure ~/.local/bin is in $PATH (add to ~/.bashrc if needed)
```

---  

## Usage

```bash
ffheic.sh -i <input_path> -o <png|jpg>
```

- `-i <input_path>` – Path to a **single HEIC file** or a **directory** containing HEIC files.  
- `-o <png|jpg>` – Desired output format. Must be either `png` or `jpg`.

The script creates (or re‑uses) a subfolder named `converted` next to the first input file and writes a timestamped log file inside that folder.

---  

## Options & Arguments

| Flag | Description | Example |
|------|-------------|---------|
| `-i` | Input file **or** directory | `-i ./photos` |
| `-o` | Output image type (`png` or `jpg`) | `-o png` |
| `-h` | Show a short help message (handled internally) | `ffheic.sh -h` |

*No other flags are required.*

---  

## Examples

### Convert an entire directory to PNG

```bash
./ffheic.sh -i /home/user/pictures/heic_collection -o png
```

- All `*.heic` files under `/home/user/pictures/heic_collection` are converted.
- Output files are placed in `/home/user/pictures/heic_collection/converted`.
- A log file like `conversion_20251123_154200.log` is created inside `converted`.

### Convert a single file to JPG

```bash
./ffheic.sh -i ./sample.heic -o jpg
```

- `sample.heic` becomes `sample.jpg` inside the same folder’s `converted` subdirectory.

### Run from Git Bash on Windows

```bash
bash ffheic.sh -i C:/Users/Me/Images -o png
```

(Use forward slashes or escape backslashes.)

---  

## How It Works

1. **Argument parsing** – validates input path and output extension.  
2. **File discovery** – uses `find` (or a single‑file check) to build an array of HEIC files.  
3. **Output folder** – creates `converted` next to the first source file; re‑uses it if it already exists.  
4. **Timestamped log** – `conversion_YYYYMMDD_HHMMSS.log` records each conversion line.  
5. **Conversion loop** – calls `ffmpeg -i source.heic dest.png|jpg` for each file, suppressing non‑error output.  

---  

## Troubleshooting

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| `ffmpeg: command not found` | ffmpeg not installed or not in `PATH` | Install ffmpeg and ensure it’s on the system `PATH`. |
| “No HEIC files found to convert.” | Wrong input directory or missing `.heic` files | Verify the path and file extensions (case‑insensitive). |
| Converted files are empty or corrupted | Out‑dated ffmpeg version | Upgrade to a recent ffmpeg release (≥ 4.0). |
| Permission denied when running script | Script not executable | Run `chmod +x ffheic.sh` again. |

---  

### License
This script is released under the MIT License – feel free to modify and redistribute.  

---  

**Happy converting!**