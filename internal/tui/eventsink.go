package tui

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/itszeeshan/subdomainx/v2/internal/types"
)

// EventSink is the interface the pipeline uses to report progress.
// It decouples the scan pipeline from the presentation layer.
type EventSink interface {
	StageStarted(stage, message string)
	StageCompleted(stage, message string)
	ToolProgress(tool, domain, status string, found int, err error)
	SubdomainsFound(results []types.SubdomainResult, totalUnique int)
	HTTPResults(results []types.HTTPResult, total int)
	PortResults(results []types.PortResult, total int)
	TakeoverResults(results []types.TakeoverResult)
	Log(level, message string)
	ScanComplete(err error)
}

// --- TUI implementation: sends tea.Msg via p.Send() ---

// TUIEventSink routes pipeline events to a Bubble Tea program.
type TUIEventSink struct {
	program *tea.Program
}

// NewTUIEventSink creates an EventSink that sends messages to the given tea.Program.
func NewTUIEventSink(p *tea.Program) *TUIEventSink {
	return &TUIEventSink{program: p}
}

func (s *TUIEventSink) StageStarted(stage, message string) {
	s.program.Send(StageMsg{Stage: stage, Status: "started", Message: message})
}

func (s *TUIEventSink) StageCompleted(stage, message string) {
	s.program.Send(StageMsg{Stage: stage, Status: "completed", Message: message})
}

func (s *TUIEventSink) ToolProgress(tool, domain, status string, found int, err error) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	s.program.Send(ToolProgressMsg{Tool: tool, Domain: domain, Status: status, Found: found, Error: errStr})
}

func (s *TUIEventSink) SubdomainsFound(results []types.SubdomainResult, totalUnique int) {
	s.program.Send(ResultMsg{Results: results, Total: totalUnique})
}

func (s *TUIEventSink) HTTPResults(results []types.HTTPResult, total int) {
	s.program.Send(HTTPResultMsg{Results: results, Total: total})
}

func (s *TUIEventSink) PortResults(results []types.PortResult, total int) {
	s.program.Send(PortResultMsg{Results: results, Total: total})
}

func (s *TUIEventSink) TakeoverResults(results []types.TakeoverResult) {
	s.program.Send(TakeoverResultMsg{Results: results})
}

func (s *TUIEventSink) Log(level, message string) {
	s.program.Send(LogMsg{Level: level, Message: message, Time: time.Now()})
}

func (s *TUIEventSink) ScanComplete(err error) {
	s.program.Send(ScanCompleteMsg{Error: err})
}

// --- CLI implementation: prints to stdout as before ---

// CLIEventSink prints pipeline events to stdout, preserving the original CLI behavior.
type CLIEventSink struct{}

// NewCLIEventSink creates an EventSink that prints to stdout.
func NewCLIEventSink() *CLIEventSink {
	return &CLIEventSink{}
}

func (s *CLIEventSink) StageStarted(stage, message string) {
	log.Println(message)
}

func (s *CLIEventSink) StageCompleted(stage, message string) {
	log.Println(message)
}

func (s *CLIEventSink) ToolProgress(tool, domain, status string, found int, err error) {
	switch status {
	case "failed":
		fmt.Printf("Error with %s on %s: %v\n", tool, domain, err)
	case "completed":
		fmt.Printf("%s found %d subdomains for %s\n", tool, found, domain)
	case "skipped":
		fmt.Printf("Skipping %s: tool not found in PATH\n", tool)
	}
}

func (s *CLIEventSink) SubdomainsFound(results []types.SubdomainResult, totalUnique int) {
	fmt.Printf("Total unique subdomains found: %d\n", totalUnique)
}

func (s *CLIEventSink) HTTPResults(_ []types.HTTPResult, total int) {
	log.Printf("HTTP scanning completed: %d results", total)
}

func (s *CLIEventSink) PortResults(_ []types.PortResult, total int) {
	log.Printf("Port scanning completed: %d results", total)
}

func (s *CLIEventSink) TakeoverResults(results []types.TakeoverResult) {
	log.Printf("Takeover check completed: %d potential vulnerabilities", len(results))
}

func (s *CLIEventSink) Log(level, message string) {
	switch level {
	case "error":
		log.Printf("ERROR: %s", message)
	case "warn":
		log.Printf("Warning: %s", message)
	default:
		log.Println(message)
	}
}

func (s *CLIEventSink) ScanComplete(err error) {
	if err != nil {
		log.Printf("Scan failed: %v", err)
	}
}
