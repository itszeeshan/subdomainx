package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadLines(t *testing.T) {
	// Create a temporary test file
	testContent := `# This is a comment
domain1.com
domain2.com

domain3.com
# Another comment
domain4.com`

	tmpFile, err := os.CreateTemp("", "test_domains")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test content
	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}
	tmpFile.Close()

	// Test reading lines
	lines, err := ReadLines(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadLines failed: %v", err)
	}

	expected := []string{"domain1.com", "domain2.com", "domain3.com", "domain4.com"}
	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, line := range expected {
		if lines[i] != line {
			t.Errorf("Expected line %d to be %s, got %s", i, line, lines[i])
		}
	}
}

func TestReadLinesEmptyFile(t *testing.T) {
	// Create an empty temporary file
	tmpFile, err := os.CreateTemp("", "test_empty")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test reading empty file
	lines, err := ReadLines(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadLines failed: %v", err)
	}

	if len(lines) != 0 {
		t.Errorf("Expected 0 lines, got %d", len(lines))
	}
}

func TestReadLinesNonExistentFile(t *testing.T) {
	// Test reading non-existent file
	_, err := ReadLines("non_existent_file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestWriteLines(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "test_output")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testLines := []string{"line1", "line2", "line3"}
	outputFile := filepath.Join(tmpDir, "subdir", "test_output.txt")

	// Test writing lines
	err = WriteLines(outputFile, testLines)
	if err != nil {
		t.Fatalf("WriteLines failed: %v", err)
	}

	// Verify the file was created and contains the expected content
	lines, err := ReadLines(outputFile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if len(lines) != len(testLines) {
		t.Errorf("Expected %d lines, got %d", len(testLines), len(lines))
	}

	for i, line := range testLines {
		if lines[i] != line {
			t.Errorf("Expected line %d to be %s, got %s", i, line, lines[i])
		}
	}
}

func TestWriteLinesEmptySlice(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "test_output")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	outputFile := filepath.Join(tmpDir, "empty_output.txt")

	// Test writing empty slice
	err = WriteLines(outputFile, []string{})
	if err != nil {
		t.Fatalf("WriteLines failed: %v", err)
	}

	// Verify the file was created
	if !FileExists(outputFile) {
		t.Error("Expected file to be created")
	}
}

func TestFileExists(t *testing.T) {
	// Test with existing file
	tmpFile, err := os.CreateTemp("", "test_exists")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	if !FileExists(tmpFile.Name()) {
		t.Error("FileExists should return true for existing file")
	}

	// Test with non-existing file
	if FileExists("non_existent_file.txt") {
		t.Error("FileExists should return false for non-existing file")
	}
}

func TestEnsureDirectory(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "test_ensure")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test creating a new subdirectory
	newDir := filepath.Join(tmpDir, "subdir1", "subdir2")
	err = EnsureDirectory(newDir)
	if err != nil {
		t.Fatalf("EnsureDirectory failed: %v", err)
	}

	// Verify directory was created
	if !FileExists(newDir) {
		t.Error("Expected directory to be created")
	}

	// Test creating an existing directory (should not fail)
	err = EnsureDirectory(newDir)
	if err != nil {
		t.Fatalf("EnsureDirectory failed for existing directory: %v", err)
	}
}

func TestEnsureDirectoryPermissionError(t *testing.T) {
	// Test with a path that would cause permission issues
	// This is a best-effort test as it depends on the system
	invalidPath := "/root/invalid/path"
	err := EnsureDirectory(invalidPath)
	if err == nil {
		// On some systems this might succeed, which is fine
		t.Log("EnsureDirectory succeeded for invalid path (system dependent)")
	}
}
