package main

import (
	"flag"
	"fmt"
	"kustomize-overlazy/kustomize"
	"log"
	"path/filepath"
)

func main() {
	var overlayAbsPath string
	var err error

	// Define a string flag to accept the overlay path
	overlayPath := flag.String("overlay", "", "Path to the Kustomize overlay")

	// Parse the flags
	flag.Parse()

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
