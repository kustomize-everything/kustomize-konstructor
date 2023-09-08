package kustomize

import (
	"io/ioutil"
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
			err := RenderSingleOverlay(tt.inputPath)
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
