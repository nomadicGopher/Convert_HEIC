// Package main provides a command-line tool for converting HEIC/HEIF images to PNG or JPEG using ffmpeg.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	outType = flag.String("output", "", "png, jpg or jpeg")
	inPath  = flag.String("input", "", "File or directory path to convert")
)

// main is the entry point for the ffheic command-line tool.
// It parses flags, validates input, checks requirements, and processes files.
func main() {
	flag.Parse()

	inPathInfo, err := validateFlags()
	checkError(err)

	err = verifyRequirements()
	checkError(err)

	err = processFiles(inPathInfo)
	checkError(err)

	fmt.Println("Processing completed successfully.")
}

// validateFlags checks the command-line flags for validity and returns information about the input path.
// It ensures the output type is supported and the input path exists.
func validateFlags() (inPathInfo os.FileInfo, err error) {
	switch *outType {
	case "jpeg", "jpg", "png":
		fmt.Println("Output Type: ", *outType)
	default:
		panic("Invalid output type. Use 'png', 'jpg' or 'jpeg'.")
	}

	if inPathInfo, err = os.Stat(*inPath); err != nil {
		return nil, err
	}

	if *inPath, err = filepath.Abs(*inPath); err != nil {
		return nil, err
	}
	fmt.Println("Input Path: ", *inPath)

	return inPathInfo, nil
}

// verifyRequirements checks that the operating system is supported and that ffmpeg with HEIC/HEIF support is installed.
func verifyRequirements() (err error) {
	osType := runtime.GOOS
	switch osType {
	case "linux", "windows", "darwin":
		// Supported
	default:
		return fmt.Errorf("Unsupported OS: %s. Only Linux, Windows, and macOS are supported.", osType)
	}

	// Verify ffmpeg is installed
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return err
	}

	// Verify ffmpeg supports heic/heif
	cmd := exec.Command("ffmpeg", "-codecs")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	if !strings.Contains(string(output), "heif") {
		return fmt.Errorf("Your ffmpeg does not support HEIC/HEIF. Please install an ffmpeg build with HEIC/HEIF support. Refer to your OS documentation or ffmpeg.org for guidance.")
	}

	fmt.Println("ffmpeg with HEIC support is installed and ready.")
	return nil
}

// processFiles converts the input file or all files in the input directory to the specified output format using ffmpeg.
// It handles both single file and directory input.
func processFiles(inPathInfo os.FileInfo) (err error) {
	// If inPath is a directory, process all files inside; otherwise, process the single file
	var inFiles []string
	if inPathInfo.IsDir() {
		entries, err := os.ReadDir(*inPath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				inFiles = append(inFiles, filepath.Join(*inPath, entry.Name()))
			}
		}
	} else {
		inFiles = append(inFiles, *inPath)
	}

	// For each file, run ffmpeg conversion
	for _, inFile := range inFiles {
		outFile := inFile + "." + *outType
		cmd := exec.Command("ffmpeg", "-i", inFile, outFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Failed to convert %s: %v\n", inFile, err)
		}
		fmt.Printf("Converted %s to %s OK.\n", inFile, outFile)
	}

	return nil
}

// checkError panics if the provided error is non-nil.
// It is used for simple error handling throughout the program.
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
