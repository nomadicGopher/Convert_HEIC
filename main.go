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

func verifyRequirements() (err error) {
	// Determine OS type (Windows, MacOS, Linux)
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
		return fmt.Errorf("Failed to run ffmpeg. Please ensure ffmpeg is installed and working: %v", err)
	}
	if !strings.Contains(string(output), "heif") {
		return fmt.Errorf("Your ffmpeg does not support HEIC/HEIF. Please install an ffmpeg build with HEIC/HEIF support. Refer to your OS documentation or ffmpeg.org for guidance: %v", err)
	}

	fmt.Println("ffmpeg with HEIC support is installed and ready.")
	return nil
}

func processFiles(inPathInfo os.FileInfo) (err error) {
	// If inPath is a directory, process all files inside; otherwise, process the single file
	var files []string
	if inPathInfo.IsDir() {
		entries, err := os.ReadDir(*inPath)
		checkError(err)
		for _, entry := range entries {
			if !entry.IsDir() {
				files = append(files, filepath.Join(*inPath, entry.Name()))
			}
		}
	} else {
		files = append(files, *inPath)
	}

	// For each file, run ffmpeg conversion
	for _, file := range files {
		outFile := file + "." + *outType
		cmd := exec.Command("ffmpeg", "-i", file, outFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Printf("Converting %s to %s...\n", file, outFile)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to convert %s: %v\n", file, err)
		}
	}

	return nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
