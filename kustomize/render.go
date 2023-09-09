package kustomize

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
)

func RenderSingleOverlay(logger *slog.Logger, baseDir string, overlayPath string, outputPath string) error {
	logger.Info("Rendering overlay: " + overlayPath)

	options := krusty.MakeDefaultOptions()
	options.PluginConfig = &types.PluginConfig{
		HelmConfig: types.HelmConfig{
			Enabled: true,
			Command: "helm",
		},
	}

	k := krusty.MakeKustomizer(options)
  fs := filesys.MakeFsOnDisk()
	m, err := k.Run(fs, overlayPath)
	if err != nil {
		log.Fatalf("Failed to render overlay: %v", err)
		return err
	}

	outputOverlayPath := outputPath + "/" + KebabOverlayPath(overlayPath)

	// Create the output directory if it does not exist
	err = os.MkdirAll(outputOverlayPath, os.ModePerm)
	if err != nil {
		return err
	}

	writer := MakeWriter(fs)
	// err = writeOutput(m, outputPath)
	err = writer.WriteIndividualFiles(outputOverlayPath, m)
	if err != nil {
		log.Fatalf("Failed to write output: %v", err)
		return err
	}

	return nil
}

func writeOutput(m resmap.ResMap, outputFilename string) error {
	// Create parent directory if it does not exist
	err := os.MkdirAll(filepath.Dir(outputFilename), os.ModePerm)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	yaml, err := m.AsYaml()
	if err != nil {
		return err
	}

	_, err = outputFile.Write(yaml)
	if err != nil {
		return err
	}

	return nil
}

func RenderOverlaysInDirectory(logger *slog.Logger, baseDir string, pattern string, outputDir string) error {
	logger.Info("Rendering overlays in directory: " + baseDir)

	// Automatically add the variants of kustomization.yaml to the end of the regex
	pattern = pattern + "/?kustomization.ya?ml"
	var matcher = regexp.MustCompile(pattern)
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		logger.Debug("Walking path: " + path)
		if matcher.MatchString(path) {
			logger.Debug("Found match: " + path)
		} else {
			logger.Debug("No match: " + path)
		}

		// If the current file matches the regex, render it
		if !info.IsDir() && matcher.MatchString(path) {
			// Remove the kustomization.ya?ml from the relative path
			path = strings.ReplaceAll(path, "/kustomization.yaml", "")
			path = strings.ReplaceAll(path, "/kustomization.yml", "")

			// TODO Setup logging with levels
			// log.Printf("Rendering overlay: %s", path)

			// Create the output directory if it does not exist
			err = os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				return err
			}

			// Determine relative path from basedir for output file name
			relOverlayPath, err := filepath.Rel(baseDir, path)

			// Call the rendering function (modified to accept output file name)
			err = RenderSingleOverlay(logger, relOverlayPath, path, outputDir)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
