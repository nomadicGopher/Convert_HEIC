#!/bin/bash
# ffheic.sh - Convert HEIC images to PNG or JPG
# -------------------------------------------------
# Usage example:
#   ./ffheic.sh -i /path/to/input -o png
#   -i  Path to a single HEIC file or a directory containing HEIC files
#   -o  Desired output format (png or jpg)
# -------------------------------------------------

set -euo pipefail   # abort on error, undefined vars, pipeline failures

# -------------------------------------------------
# 1. Ensure the script is running on a Debian-based distro
# -------------------------------------------------
if ! command -v apt >/dev/null 2>&1; then
    echo "Error: This script is intended for Debian-based systems (apt)." >&2
    exit 1
fi

# -------------------------------------------------
# Helper: print usage information and exit
# -------------------------------------------------
print_usage() {
    echo "Usage: $0 -i <input_path> -o <png|jpg>" >&2
    exit 1
}

# -------------------------------------------------
# Parse command-line options
# -------------------------------------------------
while getopts ":i:o:" opt; do
    case $opt in
        i) INPUT_PATH="$OPTARG" ;;
        o) OUTPUT_EXTENSION="${OPTARG,,}" ;;   # force lowercase
        *) print_usage ;;
    esac
done

# -------------------------------------------------
# Validate arguments
# -------------------------------------------------
[[ -z "${INPUT_PATH:-}" || -z "${OUTPUT_EXTENSION:-}" ]] && print_usage
if [[ "$OUTPUT_EXTENSION" != "png" && "$OUTPUT_EXTENSION" != "jpg" ]]; then
    echo "Error: output type must be 'png' or 'jpg'." >&2
    exit 1
fi

# -------------------------------------------------
# 2. Ensure ffmpeg (with HEIC support) is installed
# -------------------------------------------------
install_ffmpeg() {
    echo "Installing ffmpeg with HEIC support..."

    # Try the default repositories first
    sudo apt update -qq
    sudo DEBIAN_FRONTEND=noninteractive apt install -y ffmpeg libheif-dev >/dev/null 2>&1 && return

    # If libheif-dev wasn't available, add a PPA that provides newer ffmpeg
    echo "libheif-dev not found - adding PPA..."
    sudo add-apt-repository -y ppa:savoury1/ffmpeg4 >/dev/null 2>&1
    sudo apt update -qq
    sudo DEBIAN_FRONTEND=noninteractive apt install -y ffmpeg libheif-dev >/dev/null 2>&1
}

# Verify ffmpeg presence
if ! command -v ffmpeg >/dev/null 2>&1; then
    install_ffmpeg
else
    if ! ffmpeg -codecs 2>/dev/null | grep -qE 'heif'; then
        echo "ffmpeg is installed but lacks HEIC support."
        install_ffmpeg
    fi
fi

# -------------------------------------------------
# Build an array of HEIC files to process
# -------------------------------------------------
if [[ -d "$INPUT_PATH" ]]; then
    mapfile -t HEIC_FILES < <(find "$INPUT_PATH" -type f -iname "*.heic")
elif [[ -f "$INPUT_PATH" ]]; then
    HEIC_FILES=("$INPUT_PATH")
else
    echo "Error: '$INPUT_PATH' is neither a file nor a directory." >&2
    exit 1
fi

if [[ ${#HEIC_FILES[@]} -eq 0 ]]; then
    echo "No HEIC files found to convert."
    exit 0
fi

# -------------------------------------------------
# Prepare output folder (reuse if it already exists)
# -------------------------------------------------
BASE_DIRECTORY=$(dirname "${HEIC_FILES[0]}")
OUTPUT_DIRECTORY="${BASE_DIRECTORY}/converted"
mkdir -p "$OUTPUT_DIRECTORY"

# -------------------------------------------------
# Create a timestamped log file inside the output folder
# -------------------------------------------------
CURRENT_TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
LOG_FILE_PATH="${OUTPUT_DIRECTORY}/conversion_${CURRENT_TIMESTAMP}.log"

# -------------------------------------------------
# Redirect *all* output (stdout + stderr) to the log file
# -------------------------------------------------
exec >"$LOG_FILE_PATH" 2>&1

# -------------------------------------------------
# Conversion loop - everything now goes to the log
# -------------------------------------------------
echo "=== Conversion started: $(date) ==="
for source_file in "${HEIC_FILES[@]}"; do
    base_name=$(basename "$source_file" .heic)
    destination_file="${OUTPUT_DIRECTORY}/${base_name}.${OUTPUT_EXTENSION}"

    if [[ "$OUTPUT_EXTENSION" == "png" ]]; then
        # PNG: let ffmpeg pick the PNG encoder
        ffmpeg -hide_banner -loglevel error -y -i "$source_file" "$destination_file"
    else
        # JPG: force libx264 and set a compatible pixel format
        ffmpeg -hide_banner -loglevel error -y -i "$source_file" -c:v libx264 -pix_fmt yuv420p "$destination_file"
    fi

    echo "Converted: $source_file → $destination_file"
done
echo "=== Conversion finished: $(date) ==="

# -------------------------------------------------
# Inform the user (to the terminal) where the log is
# -------------------------------------------------
# `>&2` writes to the original stderr, which is still the terminal
echo "All conversions logged to: $LOG_FILE_PATH" >&2
