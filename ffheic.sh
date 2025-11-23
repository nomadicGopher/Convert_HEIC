#!/usr/bin/env bash
# ffheic.sh – Convert HEIC images to PNG or JPG
# -------------------------------------------------
# Usage example:
#   ./ffheic.sh -i /path/to/input -o png
#   -i  Path to a single HEIC file or a directory containing HEIC files
#   -o  Desired output format (png or jpg)
# -------------------------------------------------

# Exit immediately if a command exits with a non‑zero status,
# treat unset variables as an error, and propagate errors through pipelines.
set -euo pipefail   # <-- safety flags: abort on error, undefined vars, and pipeline failures

# ----------------------------------------------------------------------
# Helper: print usage information and exit
# ----------------------------------------------------------------------
print_usage() {
    echo "Usage: $0 -i <input_path> -o <png|jpg>"
    exit 1
}

# ----------------------------------------------------------------------
# Parse command‑line options
# ----------------------------------------------------------------------
while getopts ":i:o:" opt; do
    case $opt in
        i) INPUT_PATH="$OPTARG" ;;                # path supplied by user
        o) OUTPUT_EXTENSION="${OPTARG,,}" ;;     # force lowercase (png or jpg)
        *) print_usage ;;
    esac
done

# Validate required arguments
[[ -z "${INPUT_PATH:-}" || -z "${OUTPUT_EXTENSION:-}" ]] && print_usage
if [[ "$OUTPUT_EXTENSION" != "png" && "$OUTPUT_EXTENSION" != "jpg" ]]; then
    echo "Error: output type must be 'png' or 'jpg'."
    exit 1
fi

# ----------------------------------------------------------------------
# Ensure ffmpeg is available
# ----------------------------------------------------------------------
command -v ffmpeg >/dev/null 2>&1 || {
    echo "Error: ffmpeg is not installed or not in PATH."
    exit 1
}

# ----------------------------------------------------------------------
# Build an array of HEIC files to process
# ----------------------------------------------------------------------
if [[ -d "$INPUT_PATH" ]]; then
    # Input is a directory – find all *.heic files (case‑insensitive)
    mapfile -t HEIC_FILES < <(find "$INPUT_PATH" -type f -iname "*.heic")
elif [[ -f "$INPUT_PATH" ]]; then
    # Input is a single file
    HEIC_FILES=("$INPUT_PATH")
else
    echo "Error: '$INPUT_PATH' is neither a file nor a directory."
    exit 1
fi

# Nothing to do?
if [[ ${#HEIC_FILES[@]} -eq 0 ]]; then
    echo "No HEIC files found to convert."
    exit 0
fi

# ----------------------------------------------------------------------
# Prepare output folder (reuse if it already exists)
# ----------------------------------------------------------------------
# Use the directory of the first input file as the base location
BASE_DIRECTORY=$(dirname "${HEIC_FILES[0]}")
OUTPUT_DIRECTORY="${BASE_DIRECTORY}/converted"
mkdir -p "$OUTPUT_DIRECTORY"   # creates the folder if missing; otherwise reuses it

# ----------------------------------------------------------------------
# Create a timestamped log file inside the output folder
# ----------------------------------------------------------------------
CURRENT_TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
LOG_FILE_PATH="${OUTPUT_DIRECTORY}/conversion_${CURRENT_TIMESTAMP}.log"

# ----------------------------------------------------------------------
# Conversion loop – write progress to the log file
# ----------------------------------------------------------------------
{
    echo "=== Conversion started: $(date) ==="
    for source_file in "${HEIC_FILES[@]}"; do
        # Strip the .heic extension and build the destination path
        base_name=$(basename "$source_file" .heic)
        destination_file="${OUTPUT_DIRECTORY}/${base_name}.${OUTPUT_EXTENSION}"

        # Perform conversion (quiet, only error messages)
        ffmpeg -hide_banner -loglevel error -i "$source_file" "$destination_file"

        echo "Converted: $source_file → $destination_file"
    done
    echo "=== Conversion finished: $(date) ==="
} >"$LOG_FILE_PATH"

# ----------------------------------------------------------------------
# Inform the user where the log can be found
# ----------------------------------------------------------------------
echo "All conversions logged to: $LOG_FILE_PATH"