# Convert_HEIC

CLI utility that batch‑converts **HEIC** images to **PNG** or **JPG/JPEG**
using **ImageMagick**.

## Features

- Converts HEIC images to PNG, JPG, or JPEG formats.
- Supports batch conversion of all HEIC files in a directory.
- Parallel processing with configurable worker count for faster batch conversion.
  - **Default**: 4 workers

## Requirements

- **Linux**
  - Currently Windows & MacOS are not supported
- **ImageMagick**
  - ImageMagick must support HEIC format. You can check this by running
  `convert --version` and looking for "heic" in the list of supported formats.

## Usage

```sh
Convert_HEIC_{arch} -input="{filePath|directoryPath}" -output="png|jpg|jpeg" -workers=4
```

## Example

```sh
Convert_HEIC_amd64 -input="/home/username/Pictures" -output="jpg"
```
