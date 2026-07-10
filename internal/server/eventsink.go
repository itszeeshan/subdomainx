package server

import (
	"time"

	"github.com/itszeeshan/subdomainx/internal/types"
)

// APIEventSink implements tui.EventSink and routes pipeline progress
// updates into a ScanJob's fields for REST API polling.
type APIEventSink struct {
	job *ScanJob
}

// NewAPIEventSink creates an EventSink that writes into the given ScanJob.
func NewAPIEventSink(job *ScanJob) *APIEventSink {
	return &APIEventSink{job: job}
}

func (s *APIEventSink) StageStarted(stage, message string) {
	s.job.mu.Lock()
	defer s.job.mu.Unlock()
	s.job.Progress.Stage = stage
	s.job.Progress.StageMessage = message
	s.job.Status = StatusRunning
}

func (s *APIEventSink) StageCompleted(stage, message string) {
	s.job.mu.Lock()
	defer s.job.mu.Unlock()
	s.job.Progress.Stage = stage
	s.job.Progress.StageMessage = message
}

func (s *APIEventSink) ToolProgress(tool, domain, status string, found int, err error) {
	s.job.mu.Lock()
	defer s.job.mu.Unlock()
	switch status {
	case "completed":
		s.job.Progress.ToolsCompleted++
	case "failed":
		s.job.Progress.ToolsFailed++
	}
}

func (s *APIEventSink) SubdomainsFound(results []types.SubdomainResult, totalUnique int) {
	s.job.mu.Lock()
	defer s.job.mu.Unlock()
	s.job.Progress.SubdomainsFound = totalUnique
	if s.job.Results == nil {
		s.job.Results = &ScanResults{}
	}
	s.job.Results.Subdomains = results
}

func (s *APIEventSink) HTTPResults(results []types.HTTPResult, total int) {
	s.job.mu.Lock()
	defer s.job.mu.Unlock()
	s.job.Progress.HTTPResults = total
	if s.job.Results == nil {
		s.job.Results = &ScanResults{}
	}
	s.job.Results.HTTP = results
}

func (s *APIEventSink) PortResults(results []types.PortResult, total int) {
	s.job.mu.Lock()
	defer s.job.mu.Unlock()
	s.job.Progress.PortResults = total
	if s.job.Results == nil {
		s.job.Results = &ScanResults{}
	}
	s.job.Results.Ports = results
}

func (s *APIEventSink) TakeoverResults(results []types.TakeoverResult) {
	s.job.mu.Lock()
	defer s.job.mu.Unlock()
	if s.job.Results == nil {
		s.job.Results = &ScanResults{}
	}
	s.job.Results.Takeover = results
}

func (s *APIEventSink) Log(level, message string) {
	// API mode: logs are silently consumed. Could be extended to store
	// a log buffer on the job if needed.
}

func (s *APIEventSink) ScanComplete(err error) {
	s.job.mu.Lock()
	defer s.job.mu.Unlock()
	now := time.Now()
	s.job.CompletedAt = &now
	if err != nil {
		s.job.Status = StatusFailed
		s.job.Error = err.Error()
	} else {
		s.job.Status = StatusCompleted
	}
}
