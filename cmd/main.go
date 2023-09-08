package main

import (
	"flag"
	"fmt"
	"kustomize-overlazy/kustomize"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var err error
	var overlayAbsPath string
	var logger *slog.Logger

	// Define flags
	baseDir := flag.String("baseDir", ".", "Base directory to search for kustomization files")
	debugFlag := flag.Bool("debug", false, "enable debug logging")
	outputDir := flag.String("outputDir", "output", "Output directory to write rendered overlays")
	overlayPath := flag.String("overlay", "", "Path to the Kustomize overlay")
	pattern := flag.String("pattern", "kustomization.yaml", "Pattern to match kustomization files")

	// Parse the flags
	flag.Parse()

	if *debugFlag {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	// Determine which function to call based on the provided flags
	if strings.TrimSpace(*pattern) != "" {
		// Call the RenderOverlaysInDirectory function with the provided base directory and pattern
		err = kustomize.RenderOverlaysInDirectory(logger, *baseDir, *pattern, *outputDir)
		if err != nil {
			logger.Error("Failed to render overlays in directory: ", err)
		}

		fmt.Println("Overlays rendered successfully!")
	} else {
		// Check if overlay path is provided
		if *overlayPath == "" {
			logger.Error("Please provide a valid overlay path")
		}

		// Convert the relative path to an absolute path
		overlayAbsPath, err = filepath.Abs(*overlayPath)
		if err != nil {
			logger.Error("Failed to find absolute path: ", err)
		}

		// Call the RenderSingleOverlay function with the provided overlay path
		err = kustomize.RenderSingleOverlay(logger, overlayAbsPath, "output.yaml")
		if err != nil {
			logger.Error("Failed to render overlay: ", err)
		}

		logger.Info("Overlay rendered successfully!")
	}

}
