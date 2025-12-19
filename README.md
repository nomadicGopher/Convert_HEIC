## ffheic

A tiny Bash utility that batch‑converts **HEIC** images to **PNG** or **JPG** using **ffmpeg**.  
It is designed for **Debian‑based Linux distributions** (Ubuntu, Linux Mint, Pop!_OS, etc).

> [!IMPORTANT]
> ffheic is still undergoing development & testing.
>

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
| Requirement | Why it’s needed | Debian‑based install command |
|-------------|-----------------|------------------------------|
| **Python 3.8+** | Runs the Python conversion script (`ffheic.py`) | `sudo apt install -y python3` |
| **ffmpeg** (≥ 4.0) with **HEIC** support | Performs the actual image conversion | `sudo apt install -y ffmpeg libheif-dev` |
| **add‑apt‑repository** (for the optional PPA) | Allows adding the `savoury1/ffmpeg4` PPA if `libheif-dev` isn’t in the default repos | `sudo apt install -y software-properties-common` |

Make sure `ffmpeg` is reachable from your `PATH`:

```bash
ffmpeg -version   # should print version information
```

If it isn’t installed, see the [ffmpeg download page](https://ffmpeg.org/download.html) or use the commands above.

---  

## Installation

1. **Clone the repository (or copy the Python script)**  

   ```bash
   git clone https://github.com/nomadicGopher/ffheic.git
   cd ffheic
   ```

---  

## Usage

```bash
python ffheic.py -i <input_path> -o <png|jpg>
```

- `-i <input_path>` – Path to a **single HEIC file** or a **directory** containing HEIC files.  
- `-o <png|jpg>` – Desired output format. Must be either `png` or `jpg`.

The script creates (or re‑uses) a subfolder named `converted` next to the first input file. Progress and errors are written to stdout (the console).

---  

## Examples

### Convert an entire directory to PNG

```bash
python ffheic.py -i /home/user/pictures/heic_collection -o png
```

- All `*.heic` files under `/home/user/pictures/heic_collection` are converted.  
- Output files are placed in `/home/user/pictures/heic_collection/converted`.  

### Convert a single file to JPG

```bash
python ffheic.sh -i ./sample.heic -o jpg
```

- `sample.heic` becomes `sample.jpg` inside the same folder’s `converted` subdirectory.
