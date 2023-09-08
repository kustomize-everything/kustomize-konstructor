package kustomize

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
)

func RenderSingleOverlay(overlayPath string, outputFilename ...string) error {
	outputPath := "output.yaml" // default output path
	if len(outputFilename) > 0 && outputFilename[0] != "" {
		outputPath = outputFilename[0]
	}

	options := krusty.MakeDefaultOptions()
	options.PluginConfig = &types.PluginConfig{
		HelmConfig: types.HelmConfig{
			Enabled: true,
			Command: "helm",
		},
	}

	k := krusty.MakeKustomizer(options)

	m, err := k.Run(filesys.MakeFsOnDisk(), overlayPath)
	if err != nil {
		log.Fatalf("Failed to render overlay: %v", err)
		return err
	}

	err = writeOutput(m, outputPath)
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

func RenderOverlaysInDirectory(baseDir string, pattern string, outputDir string) error {
	// TODO Setup logging with levels
	// log.Println("Rendering overlays in directory: " + baseDir)
	pattern = pattern + "/?kustomization.ya?ml"
	var matcher = regexp.MustCompile(pattern)
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Automatically add the variants of kustomization.yaml to the end of the regex
		// TODO Setup logging with levels
		// log.Println("Walking path: " + path)
		// if matcher.MatchString(path) {
		// 	log.Println("Found match: " + path)
		// } else {
		// 	log.Println("No match: " + path)
		// }

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

			// Prepare the output file name based on the relative path
			outputFileName := outputDir + "/" + strings.ReplaceAll(relOverlayPath, "/", "-") + ".yaml"

			// Call the rendering function (modified to accept output file name)
			err = RenderSingleOverlay(path, outputFileName)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
