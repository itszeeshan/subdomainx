package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ReadLines reads all lines from a file and returns them as a slice of strings
func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// WriteLines writes a slice of strings to a file, one per line
func WriteLines(filename string, lines []string) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return writer.Flush()
}

// FileExists checks if a file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// EnsureDirectory creates a directory if it doesn't exist
func EnsureDirectory(dir string) error {
	return os.MkdirAll(dir, 0755)
}
