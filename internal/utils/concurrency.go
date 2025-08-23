package utils

import (
	"context"
	"sync"
	"time"
)

// WorkerPool manages a pool of workers for concurrent processing
type WorkerPool struct {
	workers     int
	rateLimit   int
	rateLimiter *time.Ticker
	jobChan     chan func()
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers, rateLimit int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	var ticker *time.Ticker
	if rateLimit > 0 {
		ticker = time.NewTicker(time.Second / time.Duration(rateLimit))
	}

	pool := &WorkerPool{
		workers:     workers,
		rateLimit:   rateLimit,
		rateLimiter: ticker,
		jobChan:     make(chan func(), workers*2),
		ctx:         ctx,
		cancel:      cancel,
	}

	pool.start()
	return pool
}

// start initializes the worker goroutines
func (wp *WorkerPool) start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// worker is the main worker function
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case job, ok := <-wp.jobChan:
			if !ok {
				return
			}

			// Apply rate limiting if configured
			if wp.rateLimiter != nil {
				<-wp.rateLimiter.C
			}

			job()

		case <-wp.ctx.Done():
			return
		}
	}
}

// Submit adds a job to the worker pool
func (wp *WorkerPool) Submit(job func()) {
	select {
	case wp.jobChan <- job:
	case <-wp.ctx.Done():
	}
}

// Wait waits for all workers to complete
func (wp *WorkerPool) Wait() {
	close(wp.jobChan)
	wp.wg.Wait()
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	wp.cancel()
	if wp.rateLimiter != nil {
		wp.rateLimiter.Stop()
	}
}

// Semaphore provides a simple semaphore implementation
type Semaphore struct {
	sem chan struct{}
}

// NewSemaphore creates a new semaphore with the given capacity
func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{
		sem: make(chan struct{}, capacity),
	}
}

// Acquire acquires a permit from the semaphore
func (s *Semaphore) Acquire() {
	s.sem <- struct{}{}
}

// Release releases a permit back to the semaphore
func (s *Semaphore) Release() {
	<-s.sem
}

// TryAcquire attempts to acquire a permit without blocking
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	ticker *time.Ticker
	stop   chan struct{}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int) *RateLimiter {
	if rate <= 0 {
		return &RateLimiter{}
	}

	return &RateLimiter{
		ticker: time.NewTicker(time.Second / time.Duration(rate)),
		stop:   make(chan struct{}),
	}
}

// Wait waits for the next tick
func (rl *RateLimiter) Wait() {
	if rl.ticker != nil {
		<-rl.ticker.C
	}
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
	if rl.ticker != nil {
		rl.ticker.Stop()
		close(rl.stop)
	}
}
