// Package main provides a command-line tool for converting HEIC/HEIF images to PNG or JPEG using ImageMagick.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	outType       = flag.String("output", "", "Output image format: png, jpg, or jpeg (required)")
	inPath        = flag.String("input", "", "File or directory path to convert (required)")
	workers       = flag.Int("workers", 4, "Number of parallel conversions (only applies to directories)")
	validOutTypes = map[string]struct{}{
		"png":  {},
		"jpg":  {},
		"jpeg": {},
	}
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -input <file|dir> -output <png|jpg|jpeg> [-workers N]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := validateRequiredFlags(); err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}

	if err := verifyRequirements(); err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}

	inPathInfo, err := validateFlags()
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}

	if err := processFiles(inPathInfo); err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}

	fmt.Fprintln(os.Stdout, "INFO: Processing completed successfully.")
}

// validateRequiredFlags ensures required flags are provided.
func validateRequiredFlags() error {
	if strings.TrimSpace(*inPath) == "" || strings.TrimSpace(*outType) == "" {
		flag.Usage()
		return errors.New("both -input and -output flags are required")
	}
	return nil
}

// verifyRequirements checks that the operating system is supported and that ImageMagick with HEIC/HEIF support is installed.
func verifyRequirements() error {
	osType := runtime.GOOS
	switch osType {
	case "linux":
		// Verify ImageMagick is installed
		if _, err := exec.LookPath("convert"); err != nil {
			return errors.New("the 'convert' command does not exist, please ensure that ImageMagick is installed and accessible via PATH")
		}

		// Check if 'convert' supports HEIC
		output, err := exec.Command("convert", "--version").CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to run 'convert --version': %v", err)
		}
		if !strings.Contains(strings.ToLower(string(output)), "heic") {
			return errors.New("ImageMagick 'convert' does not support HEIC. Try installing libheif* and then reinstall ImageMagick")
		}
	case "windows":
		return errors.New("currently, Windows is not supported")
	case "darwin":
		return errors.New("currently, Darwin/MacOS is not supported")
	default:
		return fmt.Errorf("%s is not supported", osType)
	}

	fmt.Fprintln(os.Stdout, "INFO: OS requirements are met.")
	return nil
}

// validateFlags checks the command-line flags for validity and returns information about the input path.
func validateFlags() (os.FileInfo, error) {
	absPath, err := filepath.Abs(*inPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}
	*inPath = absPath

	inPathInfo, err := os.Stat(*inPath)
	if err != nil {
		return nil, fmt.Errorf("input path error: %v", err)
	}
	fmt.Fprintln(os.Stdout, "INFO: Input Path:", *inPath)

	outTypeLower := strings.ToLower(*outType)
	if _, ok := validOutTypes[outTypeLower]; !ok {
		return nil, errors.New("invalid output type. Use 'png', 'jpg', or 'jpeg'")
	}
	*outType = outTypeLower
	fmt.Fprintln(os.Stdout, "INFO: Output Type:", *outType)

	return inPathInfo, nil
}

// processFiles converts the input file or all files in the input directory to the specified output format using ImageMagick.
// It handles both single file and directory input, and processes directories in parallel.
func processFiles(inPathInfo os.FileInfo) error {
	if inPathInfo.IsDir() {
		return processDirectory(*inPath)
	}
	return processSingleFile(*inPath)
}

// processDirectory processes all .heic files in the directory in parallel.
func processDirectory(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	var heicFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if isHeicFile(entry.Name()) {
			heicFiles = append(heicFiles, filepath.Join(dirPath, entry.Name()))
		}
	}

	if len(heicFiles) == 0 {
		return errors.New("no HEIC files found in the directory")
	}

	// Parallel processing with worker pool
	numWorkers := *workers
	if numWorkers < 1 {
		numWorkers = 1
	}
	fileCh := make(chan string, len(heicFiles))
	errCh := make(chan error, len(heicFiles))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileCh {
				if err := processSingleFile(file); err != nil {
					errCh <- err
				}
			}
		}()
	}

	for _, file := range heicFiles {
		fileCh <- file
	}
	close(fileCh)
	wg.Wait()
	close(errCh)

	var errs []string
	for e := range errCh {
		errs = append(errs, e.Error())
	}
	if len(errs) > 0 {
		return fmt.Errorf("some files failed to convert:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}

// processSingleFile converts a single HEIC file to the specified output format.
func processSingleFile(inFile string) error {
	if !isHeicFile(inFile) {
		return fmt.Errorf("file %s does not have a .heic extension", inFile)
	}
	outFile := buildOutputFilename(inFile, *outType)
	cmd := exec.Command("convert", inFile, outFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to convert %s: %v", inFile, err)
	}
	fmt.Fprintf(os.Stdout, "INFO: Converted %s to %s.\n", inFile, outFile)
	return nil
}

// isHeicFile checks if the file has a .heic extension (case-insensitive).
func isHeicFile(filename string) bool {
	return strings.EqualFold(filepath.Ext(filename), ".heic")
}

// buildOutputFilename constructs the output filename based on the input file and output type.
func buildOutputFilename(inFile, outType string) string {
	ext := filepath.Ext(inFile)
	base := strings.TrimSuffix(inFile, ext)
	return base + "." + outType
}
