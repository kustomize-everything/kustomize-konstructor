package main

import (
	"flag"
	"fmt"
	"kustomize-overlazy/kustomize"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	var overlayAbsPath string
	var err error

	// Define flags
	overlayPath := flag.String("overlay", "", "Path to the Kustomize overlay")
	baseDir := flag.String("baseDir", ".", "Base directory to search for kustomization files")
	pattern := flag.String("pattern", "kustomization.yaml", "Pattern to match kustomization files")
	outputDir := flag.String("outputDir", "output", "Output directory to write rendered overlays")

	// Parse the flags
	flag.Parse()

	// Determine which function to call based on the provided flags
	if strings.TrimSpace(*pattern) != "" {
		// Call the RenderOverlaysInDirectory function with the provided base directory and pattern
		err = kustomize.RenderOverlaysInDirectory(*baseDir, *pattern, *outputDir)
		if err != nil {
			log.Fatalf("Failed to render overlays in directory: %v", err)
		}

		fmt.Println("Overlays rendered successfully!")
	} else {
		// Check if overlay path is provided
		if *overlayPath == "" {
			log.Fatalf("Please provide a valid overlay path")
		}

		// Convert the relative path to an absolute path
		overlayAbsPath, err = filepath.Abs(*overlayPath)
		if err != nil {
			log.Fatal(err)
		}

		// Call the RenderSingleOverlay function with the provided overlay path
		err = kustomize.RenderSingleOverlay(overlayAbsPath, "output.yaml")
		if err != nil {
			log.Fatalf("Failed to render overlay: %v", err)
		}

		fmt.Println("Overlay rendered successfully!")
	}

}
