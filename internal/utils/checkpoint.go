package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/itszeeshan/subdomainx/internal/types"
)

type Checkpoint struct {
	ScanID       string                  `json:"scan_id"`
	Timestamp    time.Time               `json:"timestamp"`
	Domain       string                  `json:"domain"`
	WildcardFile string                  `json:"wildcard_file"`
	Config       map[string]interface{}  `json:"config"`
	Progress     ProgressState           `json:"progress"`
	Subdomains   []types.SubdomainResult `json:"subdomains"`
	HTTPResults  []types.HTTPResult      `json:"http_results"`
	PortResults  []types.PortResult      `json:"port_results"`
	Completed    bool                    `json:"completed"`
	ErrorMessage string                  `json:"error_message,omitempty"`
}

type ProgressState struct {
	TotalTasks     int       `json:"total_tasks"`
	CompletedTasks int       `json:"completed_tasks"`
	StartTime      time.Time `json:"start_time"`
	LastUpdate     time.Time `json:"last_update"`
}

func SaveCheckpoint(checkpoint *Checkpoint, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create checkpoint filename
	checkpointFile := filepath.Join(outputDir, fmt.Sprintf("%s_checkpoint.json", checkpoint.ScanID))

	// Marshal checkpoint to JSON
	data, err := json.MarshalIndent(checkpoint, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checkpoint: %v", err)
	}

	// Write to file
	if err := os.WriteFile(checkpointFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write checkpoint file: %v", err)
	}

	return nil
}

func LoadCheckpoint(scanID, outputDir string) (*Checkpoint, error) {
	checkpointFile := filepath.Join(outputDir, fmt.Sprintf("%s_checkpoint.json", scanID))

	// Check if checkpoint file exists
	if _, err := os.Stat(checkpointFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("checkpoint file not found: %s", checkpointFile)
	}

	// Read checkpoint file
	data, err := os.ReadFile(checkpointFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read checkpoint file: %v", err)
	}

	// Unmarshal checkpoint
	var checkpoint Checkpoint
	if err := json.Unmarshal(data, &checkpoint); err != nil {
		return nil, fmt.Errorf("failed to unmarshal checkpoint: %v", err)
	}

	return &checkpoint, nil
}

func ListCheckpoints(outputDir string) ([]string, error) {
	// Check if output directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return []string{}, nil // No checkpoints if directory doesn't exist
	}

	// Read directory
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read output directory: %v", err)
	}

	var checkpoints []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "_checkpoint.json") {
			// Extract scan ID from filename (remove _checkpoint.json suffix)
			scanID := strings.TrimSuffix(file.Name(), "_checkpoint.json")
			checkpoints = append(checkpoints, scanID)
		}
	}

	return checkpoints, nil
}

func DeleteCheckpoint(scanID, outputDir string) error {
	checkpointFile := filepath.Join(outputDir, fmt.Sprintf("%s_checkpoint.json", scanID))

	if err := os.Remove(checkpointFile); err != nil {
		return fmt.Errorf("failed to delete checkpoint file: %v", err)
	}

	return nil
}

func CreateCheckpoint(scanID, domain, wildcardFile string, config map[string]interface{}) *Checkpoint {
	return &Checkpoint{
		ScanID:       scanID,
		Timestamp:    time.Now(),
		Domain:       domain,
		WildcardFile: wildcardFile,
		Config:       config,
		Progress: ProgressState{
			StartTime:  time.Now(),
			LastUpdate: time.Now(),
		},
		Subdomains:  []types.SubdomainResult{},
		HTTPResults: []types.HTTPResult{},
		PortResults: []types.PortResult{},
		Completed:   false,
	}
}

func (c *Checkpoint) UpdateProgress(completedTasks, totalTasks int) {
	c.Progress.CompletedTasks = completedTasks
	c.Progress.TotalTasks = totalTasks
	c.Progress.LastUpdate = time.Now()
}

func (c *Checkpoint) AddSubdomains(subdomains []types.SubdomainResult) {
	c.Subdomains = append(c.Subdomains, subdomains...)
}

func (c *Checkpoint) AddHTTPResults(results []types.HTTPResult) {
	c.HTTPResults = append(c.HTTPResults, results...)
}

func (c *Checkpoint) AddPortResults(results []types.PortResult) {
	c.PortResults = append(c.PortResults, results...)
}

func (c *Checkpoint) MarkCompleted() {
	c.Completed = true
	c.Progress.LastUpdate = time.Now()
}

func (c *Checkpoint) MarkError(errorMsg string) {
	c.ErrorMessage = errorMsg
	c.Progress.LastUpdate = time.Now()
}
