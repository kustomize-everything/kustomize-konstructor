package kustomize

import (
	"io/ioutil"
	"log/slog"
	"os"
	"testing"
)

func TestRenderSingleOverlay(t *testing.T) {
	tests := []struct {
		name               string
		inputPath          string
		outputPath         string
		expectedOutputFile string
	}{
		{
			name:               "Test Case Napping Octopus",
			inputPath:          "../tests/overlays/napping_octopus",
			outputPath:         "output.yaml",
			expectedOutputFile: "../tests/overlays/napping_octopus/expected_output.yaml",
		},
		// Add more test cases here as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			err := RenderSingleOverlay(logger, tt.inputPath, "output.yaml")
			if err != nil {
				t.Fatalf("failed to render overlay: %v", err)
			}

			// Read the generated output and expected output
			generatedOutput, err := ioutil.ReadFile(tt.outputPath)
			if err != nil {
				t.Fatalf("failed to read generated output: %v", err)
			}

			expectedOutput, err := ioutil.ReadFile(tt.expectedOutputFile)
			if err != nil {
				t.Fatalf("failed to read expected output: %v", err)
			}

			// Compare the generated output with the expected output
			if string(generatedOutput) != string(expectedOutput) {
				t.Fatalf("generated output does not match expected output")
			}

			// Cleanup: Remove the generated output file
			err = os.Remove(tt.outputPath)
			if err != nil {
				t.Fatalf("failed to remove generated output file: %v", err)
			}
		})
	}
}

func TestRenderOverlaysInDirectory(t *testing.T) {
	tests := []struct {
		name            string
		baseDir         string
		pattern         string
		expectedOutputs map[string]string // Map of output file paths to expected output file paths
	}{
		{
			name:    "Test Multiple Overlays",
			baseDir: "../tests/overlays/",
			pattern: "overlays/.*/",
			expectedOutputs: map[string]string{
				"output/napping_octopus.yaml": "../tests/overlays/napping_octopus/expected_output.yaml",
				"output/snoring_squid.yaml":   "../tests/overlays/snoring_squid/expected_output.yaml",
			},
		},
		{
			name:    "Test Multiple Overlays without trailing slash in pattern",
			baseDir: "../tests/overlays/",
			pattern: "overlays/.*",
			expectedOutputs: map[string]string{
				"output/napping_octopus.yaml": "../tests/overlays/napping_octopus/expected_output.yaml",
				"output/snoring_squid.yaml":   "../tests/overlays/snoring_squid/expected_output.yaml",
			},
		},
		{
			name:    "Test Multiple Overlays from tests root",
			baseDir: "../tests/",
			pattern: "overlays/.*/",
			expectedOutputs: map[string]string{
				"output/overlays-napping_octopus.yaml": "../tests/overlays/napping_octopus/expected_output.yaml",
				"output/overlays-snoring_squid.yaml":   "../tests/overlays/snoring_squid/expected_output.yaml",
			},
		},
		// Add more test cases here as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			err := RenderOverlaysInDirectory(logger, tt.baseDir, tt.pattern, "output")
			if err != nil {
				t.Fatalf("failed to render overlays: %v", err)
			}

			for outputPath, expectedOutputPath := range tt.expectedOutputs {
				// Read the generated output and expected output
				generatedOutput, err := ioutil.ReadFile(outputPath)
				if err != nil {
					t.Fatalf("failed to read generated output from %s: %v", outputPath, err)
				}

				expectedOutput, err := ioutil.ReadFile(expectedOutputPath)
				if err != nil {
					t.Fatalf("failed to read expected output from %s: %v", expectedOutputPath, err)
				}

				// Compare the generated output with the expected output
				if string(generatedOutput) != string(expectedOutput) {
					t.Fatalf("generated output does not match expected output for file %s", outputPath)
				}

				// Cleanup: Remove the generated output file
				err = os.Remove(outputPath)
				if err != nil {
					t.Fatalf("failed to remove generated output file %s: %v", outputPath, err)
				}
			}
		})
	}
}
