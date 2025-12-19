#!/usr/bin/env python3
"""
ffheic.py - Convert HEIC images to PNG or JPG using ffmpeg

Usage examples:
  python ffheic.py -i /path/to/input -o png
  python ffheic.py -i /path/to/file.heic -o jpg

Notes:
  - Requires `ffmpeg` to be available on PATH and compiled with
    HEIC (`libheif`) support. The script will detect ffmpeg and
    attempt to warn when HEIC codec support is missing.
  - Produces converted images in a `converted/` directory next to
    the first discovered source file.
"""

from __future__ import annotations

import argparse
from typing import Any
import shutil
import subprocess
from dataclasses import dataclass
from datetime import datetime
from pathlib import Path

@dataclass
class Config:
    """Runtime configuration from CLI arguments."""
    input_path: Path  # File or directory to process
    output_ext: str   # Output extension: 'png' or 'jpg'


def parse_args() -> Config:
    """Parse and validate command-line arguments."""

    parser = argparse.ArgumentParser(description="Convert HEIC images to PNG or JPG using ffmpeg")
    parser.add_argument("-i", "--input", required=True, help="Path to HEIC file or directory")
    parser.add_argument("-o", "--output", required=True, choices=["png", "jpg"], help="Output format: png or jpg")

    args = parser.parse_args()
    input_path = Path(args.input).expanduser().resolve()
    if not input_path.exists():
        parser.error(f"Input path does not exist: {input_path}")

    return Config(input_path=input_path, output_ext=args.output.lower())


def find_ffmpeg() -> Path | None:
    """Return the path to `ffmpeg` if available, else None."""

    ffmpeg_path = shutil.which("ffmpeg")
    return Path(ffmpeg_path) if ffmpeg_path else None


def ffmpeg_has_heif_support(ffmpeg_path: Path) -> bool:
    """Return True if ffmpeg supports HEIF/HEIC input (heuristic)."""

    try:
        proc = subprocess.run([str(ffmpeg_path), "-codecs"], capture_output=True, text=True, check=False)
        return "heif" in proc.stdout.lower() or "heic" in proc.stdout.lower()
    except Exception:
        return False


def collect_heic_files(input_path: Path) -> list[Path]:
    """Return all HEIC files under input_path (recursively if dir)."""
    if input_path.is_file():
        return [input_path]
    return [p for p in input_path.rglob("*.heic") if p.is_file()]


def prepare_output_dir(first_source: Path) -> Path:
    """Create or reuse a 'converted' output directory next to source."""
    out_dir = first_source.parent / "converted"
    out_dir.mkdir(parents=True, exist_ok=True)
    return out_dir




def log(level: str, message: str, *args: Any) -> None:
    """Print a timestamped log message to stdout."""
    ts = datetime.now().isoformat()
    if args:
        try:
            message = message % args
        except Exception:
            message = message + " " + " ".join(map(str, args))
        print(f"{ts} {level} {message}")


def convert_file(ffmpeg_path: Path, source: Path, destination: Path, output_ext: str) -> None:
    """Convert a single HEIC file to the desired output using ffmpeg.
    Raises subprocess.CalledProcessError on failure."""
    if output_ext == "png":
        cmd = [str(ffmpeg_path), "-hide_banner", "-loglevel", "error", "-y", "-i", str(source), str(destination)]
    else:
        # For JPEG output, set a pixel format compatible with most decoders
        cmd = [
            str(ffmpeg_path),
            "-hide_banner",
            "-loglevel",
            "error",
            "-y",
            "-i",
            str(source),
            "-c:v",
            "libx264",
            "-pix_fmt",
            "yuv420p",
            str(destination),
        ]
    log("INFO", "Converting: %s → %s", source, destination)
    subprocess.run(cmd, check=True)


def main() -> None:
    config = parse_args()

    ffmpeg_path = find_ffmpeg()
    if ffmpeg_path is None:
        raise SystemExit("ffmpeg not found on PATH. Install ffmpeg and try again.")

    if not ffmpeg_has_heif_support(ffmpeg_path):
        log("WARNING", "ffmpeg was found but may lack HEIC/HEIF support.")

    heic_files = collect_heic_files(config.input_path)
    if not heic_files:
        print("No HEIC files found to convert.")
        return

    out_dir = prepare_output_dir(heic_files[0])
    log("INFO", "=== Conversion started: %s ===", datetime.now().isoformat())

    for file in heic_files:
        dest_name = file.stem + "." + config.output_ext
        destination = out_dir / dest_name
        try:
            convert_file(ffmpeg_path, file, destination, config.output_ext)
        except subprocess.CalledProcessError as exc:
            log("ERROR", "Failed to convert %s: %s", file, exc)
        else:
            log("INFO", "Converted: %s → %s", file, destination)

    log("INFO", "=== Conversion finished: %s ===", datetime.now().isoformat())


if __name__ == "__main__":
    main()
