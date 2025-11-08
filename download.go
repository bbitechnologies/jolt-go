package jolt

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
	githubRepo = "bbitechnologies/jolt-go"
	releaseTag = "v0.1.0" // Update this when creating new releases
	libDir     = "lib"
)

func init() {
	// Automatically download binaries if they don't exist
	if err := ensureLibrariesExist(); err != nil {
		// Non-fatal: let CGO fail with a clearer error if binaries are missing
		fmt.Fprintf(os.Stderr, "Warning: Failed to download libraries: %v\n", err)
		fmt.Fprintf(os.Stderr, "You may need to manually download binaries from https://github.com/%s/releases\n", githubRepo)
	}
}

// ensureLibrariesExist checks if the required libraries exist and downloads them if needed
func ensureLibrariesExist() error {
	platform := runtime.GOOS + "_" + runtime.GOARCH
	platformDir := filepath.Join(libDir, platform)

	// Define required libraries for this platform
	requiredLibs := []string{
		"libJolt.a",
		"libjolt_wrapper.a",
	}

	// Check if all libraries exist
	allExist := true
	for _, lib := range requiredLibs {
		libPath := filepath.Join(platformDir, lib)
		if _, err := os.Stat(libPath); os.IsNotExist(err) {
			allExist = false
			break
		}
	}

	if allExist {
		return nil // All libraries present, nothing to do
	}

	// Libraries missing, download them
	fmt.Printf("Downloading pre-built Jolt Physics binaries for %s/%s...\n", runtime.GOOS, runtime.GOARCH)

	// Create platform directory if it doesn't exist
	if err := os.MkdirAll(platformDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", platformDir, err)
	}

	// Download each library
	for _, lib := range requiredLibs {
		libPath := filepath.Join(platformDir, lib)

		// Skip if already exists
		if _, err := os.Stat(libPath); err == nil {
			fmt.Printf("  %s already exists, skipping\n", lib)
			continue
		}

		// Construct download URL
		// Format: https://github.com/owner/repo/releases/download/tag/platform_filename
		filename := platform + "_" + lib
		url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", githubRepo, releaseTag, filename)

		fmt.Printf("  Downloading %s...\n", lib)
		if err := downloadFile(url, libPath); err != nil {
			return fmt.Errorf("failed to download %s: %w", lib, err)
		}
		fmt.Printf("  ✓ %s downloaded successfully\n", lib)
	}

	fmt.Println("All binaries downloaded successfully!")
	return nil
}

// downloadFile downloads a file from the given URL to the specified path
func downloadFile(url, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
