package kustomize

import (
	"log"
	"os"

	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
)

func RenderSingleOverlay(overlayPath string) error {
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

	err = writeOutput(m)
	if err != nil {
		log.Fatalf("Failed to write output: %v", err)
		return err
	}

	return nil
}

func writeOutput(m resmap.ResMap) error {
	outputFile, err := os.Create("output.yaml")
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

func RenderBulkOverlays(pattern string) error {
	// Crawl the current directory to find Kustomization files matching the pattern
	// For each overlay found, call the RenderSingleOverlay function to render it

	return nil
}
