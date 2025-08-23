package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type SignalHandler struct {
	interrupted bool
	checkpoint  *Checkpoint
	outputDir   string
}

func NewSignalHandler(checkpoint *Checkpoint, outputDir string) *SignalHandler {
	return &SignalHandler{
		interrupted: false,
		checkpoint:  checkpoint,
		outputDir:   outputDir,
	}
}

func (sh *SignalHandler) Start() {
	// Create channel to receive OS signals
	sigChan := make(chan os.Signal, 1)

	// Register signals to listen for
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start goroutine to handle signals
	go func() {
		sig := <-sigChan
		fmt.Printf("\n\n⚠️  Received signal %v. Saving checkpoint and shutting down gracefully...\n", sig)

		sh.interrupted = true

		// Save checkpoint before exiting
		if sh.checkpoint != nil {
			sh.checkpoint.MarkError("Scan interrupted by user")
			if err := SaveCheckpoint(sh.checkpoint, sh.outputDir); err != nil {
				fmt.Printf("❌ Failed to save checkpoint: %v\n", err)
			} else {
				fmt.Printf("✅ Checkpoint saved: %s_checkpoint.json\n", sh.checkpoint.ScanID)
				fmt.Printf("💡 Resume with: subdomainx --resume %s\n", sh.checkpoint.ScanID)
			}
		}

		os.Exit(1)
	}()
}

func (sh *SignalHandler) IsInterrupted() bool {
	return sh.interrupted
}

func (sh *SignalHandler) UpdateCheckpoint(checkpoint *Checkpoint) {
	sh.checkpoint = checkpoint
}
