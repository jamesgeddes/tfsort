package tsort_test

import (
	"os"
	"testing"

	"github.com/AlexNabokikh/tfsort/tsort"
)

const (
	validFilePath = "testdata/valid.tf"
	validTofuPath = "testdata/valid.tofu"
	outputFile    = "output.tf"
)

func TestCanIngest(t *testing.T) {
	ingestor := tsort.NewIngestor()

	t.Run("Valid Terraform File", func(t *testing.T) {
		if err := ingestor.CanIngest(validFilePath); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Valid OpenTofu File", func(t *testing.T) {
		if err := ingestor.CanIngest(validTofuPath); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("File not exists", func(t *testing.T) {
		if err := ingestor.CanIngest("notExistFile.tf"); err == nil {
			t.Errorf("Expected error but not occurred")
		}
	})

	t.Run("Invalid File Type", func(t *testing.T) {
		if err := os.WriteFile("invalid_file.txt", []byte("data"), 0o600); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := ingestor.CanIngest("invalid_file.txt"); err == nil {
			t.Errorf("Expected error but not occurred")
		}
	})

	t.Run("File with read error", func(t *testing.T) {
		if err := os.WriteFile("unreadable_file.tf", []byte("data"), 0o000); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := ingestor.CanIngest("unreadable_file.tf"); err == nil {
			t.Errorf("Expected error but not occurred")
		}
	})

	t.Run("Invalid File block", func(t *testing.T) {
		if err := os.WriteFile("invalid_file.tf", []byte("data"), 0o600); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := ingestor.CanIngest("invalid_file.tf"); err == nil {
			t.Errorf("Expected error but not occurred")
		}
	})

	// cleanup
	if err := os.Remove("invalid_file.tf"); err != nil && !os.IsNotExist(err) {
		t.Errorf("Failed to cleanup invalid_file.tf: %v", err)
	}
	if err := os.Remove("unreadable_file.tf"); err != nil && !os.IsNotExist(err) {
		t.Errorf("Failed to cleanup unreadable_file.tf: %v", err)
	}
	if err := os.Remove("invalid_file.txt"); err != nil && !os.IsNotExist(err) {
		t.Errorf("Failed to cleanup invalid_file.txt: %v", err)
	}
}

func TestParse(t *testing.T) {
	ingestor := tsort.NewIngestor()

	t.Run("Can't ingest", func(t *testing.T) {
		if err := ingestor.Parse("notExistFile.tf", outputFile, false); err == nil {
			t.Errorf("Expected error but not occurred")
		}
	})

	t.Run("Write to output file from .tf", func(t *testing.T) {
		if err := os.Remove(outputFile); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Failed to remove output file: %v", err)
		}

		if err := ingestor.Parse(validFilePath, outputFile, false); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Errorf("Output file not created")
		}

		outFile, _ := os.ReadFile(outputFile)
		expectedFile, _ := os.ReadFile("testdata/expected.tf")

		if string(outFile) != string(expectedFile) {
			t.Errorf("Output file content is not as expected")
		}
	})

	t.Run("Write to output file from .tofu", func(t *testing.T) {
		if err := os.Remove(outputFile); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Failed to remove output file: %v", err)
		}

		if err := ingestor.Parse(validTofuPath, outputFile, false); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Errorf("Output file not created")
		}

		outFile, _ := os.ReadFile(outputFile)
		expectedFile, _ := os.ReadFile("testdata/expected.tf")

		if string(outFile) != string(expectedFile) {
			t.Errorf("Output file content from .tofu is not as expected")
		}
	})

	t.Run("Write to stdout", func(t *testing.T) {
		if err := os.Remove(outputFile); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Failed to remove output file: %v", err)
		}

		outputPath := ""

		if err := ingestor.Parse(validFilePath, outputPath, true); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		outputFileInfo, err := os.Stat(outputFile)

		if outputFileInfo != nil || !os.IsNotExist(err) {
			t.Errorf("output file should not be created")
		}
	})

	t.Run("Error writing to output file", func(t *testing.T) {
		if err := os.Remove(outputFile); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Failed to remove output file: %v", err)
		}

		if err := os.WriteFile(outputFile, []byte("data"), 0o000); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := ingestor.Parse(validFilePath, outputFile, false); err == nil {
			t.Errorf("Expected error but not occurred")
		}
	})

	// cleanup
	if err := os.Remove(outputFile); err != nil && !os.IsNotExist(err) {
		t.Errorf("Failed to cleanup output file: %v", err)
	}
	if err := os.Remove("invalid_file.txt"); err != nil && !os.IsNotExist(err) {
		t.Errorf("Failed to cleanup invalid_file.txt: %v", err)
	}
}

func TestValidateFilePath(t *testing.T) {
	t.Run("File path is empty", func(t *testing.T) {
		if err := tsort.ValidateFilePath(""); err == nil {
			t.Errorf("Expected error but not occurred")
		}
	})

	t.Run("File not exists", func(t *testing.T) {
		if err := tsort.ValidateFilePath("notExistFile.tf"); err == nil {
			t.Errorf("Expected error but not occurred")
		}
	})

	t.Run("File is directory", func(t *testing.T) {
		if err := tsort.ValidateFilePath("testdata"); err == nil {
			t.Errorf("Expected error but not occurred")
		}
	})

	t.Run("Valid .tf File Path", func(t *testing.T) {
		if err := tsort.ValidateFilePath(validFilePath); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Valid .tofu File Path", func(t *testing.T) {
		if err := tsort.ValidateFilePath(validTofuPath); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	// cleanup
	if err := os.Remove("invalid_file.txt"); err != nil && !os.IsNotExist(err) {
		t.Errorf("Failed to cleanup invalid_file.txt: %v", err)
	}
}
