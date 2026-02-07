package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ContextFile represents the structure of the .ctx file
type ContextFile struct {
	Status    string    `json:"status"` // "BUSY", "IDLE"
	Task      string    `json:"task,omitempty"`
	StartedAt time.Time `json:"started_at,omitempty"`
}

const (
	StatusBusy = "BUSY"
	StatusIdle = "IDLE"
	CtxFileName = ".ctx"
)

func main() {
	// Subcommands
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'init', 'start', 'stop', 'status' or 'ls' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		runInit()
	case "start":
		startCmd.Parse(os.Args[2:])
		if startCmd.NArg() < 1 {
			fmt.Println("start requires a task name")
			os.Exit(1)
		}
		runStart(startCmd.Arg(0))
	case "stop":
		stopCmd.Parse(os.Args[2:])
		runStop()
	case "status":
		statusCmd.Parse(os.Args[2:])
		runStatus()
	case "ls":
		lsCmd.Parse(os.Args[2:])
		runLs()
	default:
		fmt.Println("expected 'init', 'start', 'stop', 'status' or 'ls' subcommands")
		os.Exit(1)
	}
}

func getCtxPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, CtxFileName), nil
}

func loadCtx() (*ContextFile, error) {
	path, err := getCtxPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var ctx ContextFile
	err = json.Unmarshal(data, &ctx)
	if err != nil {
		return nil, err
	}
	return &ctx, nil
}

func saveCtx(ctx *ContextFile) error {
	path, err := getCtxPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func runInit() {
	path, err := getCtxPath()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if _, err := os.Stat(path); err == nil {
		fmt.Println("Context file already exists.")
		return
	}
	ctx := &ContextFile{Status: StatusIdle}
	if err := saveCtx(ctx); err != nil {
		fmt.Printf("Error creating .ctx: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Initialized .ctx file")
}

func runStart(task string) {
	ctx, err := loadCtx()
	if err != nil {
		if os.IsNotExist(err) {
			ctx = &ContextFile{} // Create implicit if missing? Requirement says init, but let's be flexible or strict. Strict for now.
			fmt.Println("Error: .ctx file not found. Run 'init' first.")
			os.Exit(1)
		} else {
			fmt.Printf("Error loading .ctx: %v\n", err)
			os.Exit(1)
		}
	}
	ctx.Status = StatusBusy
	ctx.Task = task
	ctx.StartedAt = time.Now()
	if err := saveCtx(ctx); err != nil {
		fmt.Printf("Error saving .ctx: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Started task: %s\n", task)
}

func runStop() {
	ctx, err := loadCtx()
	if err != nil {
		fmt.Printf("Error: %v\n", err) // Likely not initialized
		return
	}
	ctx.Status = StatusIdle
	ctx.Task = ""
	if err := saveCtx(ctx); err != nil {
		fmt.Printf("Error saving .ctx: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Stopped task")
}

func runStatus() {
	ctx, err := loadCtx()
	if err != nil {
		// If no context file, we are neutral/idle effectively
		return
	}
	// Output format optimized for shell parsing or direct display
	// If BUSY, output the task.
	if ctx.Status == StatusBusy {
		fmt.Printf("BUSY: %s\n", ctx.Task)
	} else {
		fmt.Println("IDLE")
	}
}

func runLs() {
	// Simple listing for now
	ctx, err := loadCtx()
	if err != nil {
		fmt.Println("No context active")
		return
	}
	fmt.Printf("Status: %s\nTask: %s\nStarted: %s\n", ctx.Status, ctx.Task, ctx.StartedAt.Format(time.RFC3339))
}
