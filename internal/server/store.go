package server

import (
	"fmt"
	"sync"
)

// Store is a thread-safe in-memory store for scan jobs.
type Store struct {
	mu   sync.RWMutex
	jobs map[string]*ScanJob
}

// NewStore creates an empty scan store.
func NewStore() *Store {
	return &Store{jobs: make(map[string]*ScanJob)}
}

// Create adds a new scan job to the store.
func (s *Store) Create(job *ScanJob) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
}

// Get retrieves a scan job by ID. Returns nil if not found.
func (s *Store) Get(id string) *ScanJob {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.jobs[id]
}

// List returns a summary of all scan jobs.
func (s *Store) List() []ScanListItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]ScanListItem, 0, len(s.jobs))
	for _, job := range s.jobs {
		job.mu.RLock()
		subCount := 0
		if job.Results != nil {
			subCount = len(job.Results.Subdomains)
		}
		items = append(items, ScanListItem{
			ID:          job.ID,
			Status:      job.Status,
			Domain:      job.Domain,
			StartedAt:   job.StartedAt,
			CompletedAt: job.CompletedAt,
			Subdomains:  subCount,
		})
		job.mu.RUnlock()
	}
	return items
}

// Delete cancels a running scan and marks it as cancelled.
// Returns an error if the job is not found.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return fmt.Errorf("scan %s not found", id)
	}

	job.mu.Lock()
	defer job.mu.Unlock()

	if job.cancel != nil {
		job.cancel()
	}
	if job.Status == StatusRunning || job.Status == StatusQueued {
		job.Status = StatusCancelled
	}
	return nil
}
