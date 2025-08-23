package utils

import (
	"fmt"
	"sync"
	"time"
)

type ProgressTracker struct {
	total       int
	completed   int
	startTime   time.Time
	mu          sync.Mutex
	description string
}

func NewProgressTracker(total int, description string) *ProgressTracker {
	return &ProgressTracker{
		total:       total,
		completed:   0,
		startTime:   time.Now(),
		description: description,
	}
}

func (p *ProgressTracker) Increment() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.completed++
}

func (p *ProgressTracker) Update(completed int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.completed = completed
}

func (p *ProgressTracker) GetProgress() (completed, total int, percentage float64, eta time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	completed = p.completed
	total = p.total

	if total > 0 {
		percentage = float64(completed) / float64(total) * 100
	}

	// Calculate ETA
	if completed > 0 {
		elapsed := time.Since(p.startTime)
		rate := float64(completed) / elapsed.Seconds()
		if rate > 0 {
			remaining := float64(total-completed) / rate
			eta = time.Duration(remaining) * time.Second
		}
	}

	return
}

func (p *ProgressTracker) PrintProgress() {
	completed, total, percentage, eta := p.GetProgress()

	// Create progress bar
	barLength := 30
	filled := int(float64(barLength) * percentage / 100)
	bar := "["
	for i := 0; i < barLength; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "]"

	// Format ETA
	etaStr := "N/A"
	if eta > 0 {
		if eta.Hours() >= 1 {
			etaStr = fmt.Sprintf("%.0fh %dm", eta.Hours(), int(eta.Minutes())%60)
		} else if eta.Minutes() >= 1 {
			etaStr = fmt.Sprintf("%.0fm %ds", eta.Minutes(), int(eta.Seconds())%60)
		} else {
			etaStr = fmt.Sprintf("%.0fs", eta.Seconds())
		}
	}

	fmt.Printf("\r%s %s %d/%d (%.1f%%) ETA: %s",
		p.description, bar, completed, total, percentage, etaStr)
}

func (p *ProgressTracker) Finish() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.completed = p.total
	fmt.Println() // New line after progress bar
}

// Global progress tracker for enumeration
var (
	enumProgress *ProgressTracker
	enumMutex    sync.Mutex
)

func StartEnumerationProgress(total int) {
	enumMutex.Lock()
	defer enumMutex.Unlock()
	enumProgress = NewProgressTracker(total, "🔍 Enumerating")
}

func UpdateEnumerationProgress(completed int) {
	enumMutex.Lock()
	defer enumMutex.Unlock()
	if enumProgress != nil {
		enumProgress.Update(completed)
		enumProgress.PrintProgress()
	}
}

func IncrementEnumerationProgress() {
	enumMutex.Lock()
	defer enumMutex.Unlock()
	if enumProgress != nil {
		enumProgress.Increment()
		enumProgress.PrintProgress()
	}
}

func FinishEnumerationProgress() {
	enumMutex.Lock()
	defer enumMutex.Unlock()
	if enumProgress != nil {
		enumProgress.Finish()
		enumProgress = nil
	}
}
