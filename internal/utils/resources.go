package utils

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

type ResourceMonitor struct {
	startTime     time.Time
	startMem      runtime.MemStats
	peakMem       uint64
	checkInterval time.Duration
}

func NewResourceMonitor() *ResourceMonitor {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return &ResourceMonitor{
		startTime:     time.Now(),
		startMem:      mem,
		peakMem:       mem.Alloc,
		checkInterval: 5 * time.Second, // Check every 5 seconds
	}
}

func (rm *ResourceMonitor) Start() {
	fmt.Println("🔧 Resource monitoring started...")
}

func (rm *ResourceMonitor) Check() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	// Update peak memory
	if mem.Alloc > rm.peakMem {
		rm.peakMem = mem.Alloc
	}

	// Calculate memory usage
	memUsageMB := float64(mem.Alloc) / 1024 / 1024
	peakMemMB := float64(rm.peakMem) / 1024 / 1024

	// Get CPU info
	numCPU := runtime.NumCPU()
	numGoroutines := runtime.NumGoroutine()

	// Print resource status
	fmt.Printf("📊 Memory: %.1fMB (Peak: %.1fMB) | CPU: %d cores | Goroutines: %d\n",
		memUsageMB, peakMemMB, numCPU, numGoroutines)

	// Provide optimization recommendations
	if memUsageMB > 500 {
		fmt.Println("⚠️  High memory usage detected. Consider reducing thread count.")
	}

	if numGoroutines > numCPU*10 {
		fmt.Println("⚠️  High goroutine count detected. Consider reducing concurrency.")
	}
}

func (rm *ResourceMonitor) Optimize() {
	// Force garbage collection
	runtime.GC()
	debug.FreeOSMemory()

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	fmt.Printf("🧹 Memory optimization completed. Current usage: %.1fMB\n",
		float64(mem.Alloc)/1024/1024)
}

func (rm *ResourceMonitor) GetStats() map[string]interface{} {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	elapsed := time.Since(rm.startTime)

	return map[string]interface{}{
		"elapsed_time":   elapsed.String(),
		"current_memory": float64(mem.Alloc) / 1024 / 1024,
		"peak_memory":    float64(rm.peakMem) / 1024 / 1024,
		"num_cpu":        runtime.NumCPU(),
		"num_goroutines": runtime.NumGoroutine(),
		"gc_count":       mem.NumGC,
		"total_alloc":    float64(mem.TotalAlloc) / 1024 / 1024,
	}
}

func (rm *ResourceMonitor) PrintFinalStats() {
	stats := rm.GetStats()

	fmt.Println("\n📈 Resource Usage Summary:")
	fmt.Printf("⏱️  Total time: %s\n", stats["elapsed_time"])
	fmt.Printf("💾 Peak memory: %.1fMB\n", stats["peak_memory"])
	fmt.Printf("🔄 Total allocations: %.1fMB\n", stats["total_alloc"])
	fmt.Printf("🧹 Garbage collections: %d\n", stats["gc_count"])
	fmt.Printf("⚡ Final goroutines: %d\n", stats["num_goroutines"])
}

// Global resource monitor
var globalMonitor *ResourceMonitor

func StartResourceMonitoring() {
	globalMonitor = NewResourceMonitor()
	globalMonitor.Start()
}

func CheckResources() {
	if globalMonitor != nil {
		globalMonitor.Check()
	}
}

func OptimizeResources() {
	if globalMonitor != nil {
		globalMonitor.Optimize()
	}
}

func StopResourceMonitoring() {
	if globalMonitor != nil {
		globalMonitor.PrintFinalStats()
		globalMonitor = nil
	}
}
