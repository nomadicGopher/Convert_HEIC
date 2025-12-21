#!/usr/bin/env python3

# import subprocess
import argparse
from pathlib import Path


def parse_args() -> dict:
    """Parse and validate command-line arguments."""

    parser = argparse.ArgumentParser(
        description="Convert HEIC images to PNG or JPG using ffmpeg"
    )

    parser.add_argument(
        "-i", "--input", required=True, help="Path to HEIC file or directory"
    )

    parser.add_argument(
        "-o",
        "--output",
        required=True,
        choices=["png", "jpg"],
        help="Output format: png or jpg",
    )

    args = parser.parse_args()
    input_path = Path(args.input).expanduser().resolve()
    if not input_path.exists():
        parser.error(f"Input path does not exist: {input_path}")

    return {"input": input_path, "output": args.output.lower()}


def main() -> None:
    args = parse_args()  # args['input']


if __name__ == "__main__":
    main()
