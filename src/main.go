package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	detectCmd := flag.NewFlagSet("detect", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
	signalCmd := flag.NewFlagSet("signal", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'detect', 'list', 'summary', or 'signal' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "detect":
		detectCmd.Parse(os.Args[2:])
		if detectCmd.NArg() < 1 {
			fmt.Println("detect requires a tty path")
			os.Exit(1)
		}
		runDetect(detectCmd.Arg(0))
	case "list":
		showIdle := listCmd.Bool("show-idle", false, "Show idle panes as well")
		listCmd.Parse(os.Args[2:])
		runList(*showIdle)
	case "summary":
		summaryCmd.Parse(os.Args[2:])
		runSummary()
	case "signal":
		status := signalCmd.String("status", "", "Status: running|waiting")
		agent := signalCmd.String("agent", "unknown", "Agent name")
		signalCmd.Parse(os.Args[2:])
		if *status == "" {
			fmt.Println("signal requires --status <running|waiting>")
			os.Exit(1)
		}
		runSignal(*status, *agent)
	default:
		fmt.Println("expected 'detect', 'list', 'summary', or 'signal' subcommands")
		os.Exit(1)
	}
}

func runDetect(tty string) {
	// Original behavior: detect from ps
	// Optional: we could update the central state here too?
	// For now, keep it pure as a detector.
	state := DetectState(tty)
	fmt.Println(state)
}

func runSignal(status, agent string) {
	tty, err := getSelfTTY()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting TTY: %v\n", err)
		os.Exit(1)
	}

	// Normalize status
	status = strings.ToLower(status)
	if status != "running" && status != "waiting" {
		// Allow synonyms if needed, or strict?
		// "busy" -> "running"?
		if status == "busy" {
			status = "running"
		}
		if status == "idle" {
			status = "waiting"
		}
	}

	err = UpdateAgentState(tty, status, agent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating state: %v\n", err)
		os.Exit(1)
	}
}

func runSummary() {
	// 1. Get real active panes to prune stale state
	panes, err := GetPanes()
	if err != nil {
		// If tmux isn't running or error, just return empty
		return
	}

	activeTTYs := make([]string, 0, len(panes))
	for _, p := range panes {
		activeTTYs = append(activeTTYs, p.TTY)
	}

	// 2. Prune state
	err = PruneState(activeTTYs)
	if err != nil {
		// Ignore error, continue with best effort
	}

	// 3. Read state
	state, err := LoadState()
	if err != nil {
		return
	}

	totalAI := 0
	busyAI := 0

	// We only count agents that have registered themselves via signal
	// OR we could try to auto-detect mixed usage.
	// Requirement: "Overall Status Counter... [busy AI: 2/5]"
	// If we rely purely on signals, we iterate the state map.

	for _, s := range state.Panes {
		totalAI++
		if s.Status == "running" {
			busyAI++
		}
	}

	if totalAI > 0 {
		fmt.Printf("[busy AI: %d/%d]", busyAI, totalAI)
	}
}

func runList(showIdle bool) {
	panes, err := GetPanes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing panes: %v\n", err)
		os.Exit(1)
	}

	// Load centralized state to augment list
	state, _ := LoadState()

	for _, p := range panes {
		branch := GetGitBranch(p.Path)

		// Priority: Signal State > Detect State
		var statusMarker, desc string

		if s, ok := state.Panes[p.TTY]; ok {
			desc = s.Agent
			if s.Status == "running" {
				statusMarker = "BUSY"
			} else {
				statusMarker = "IDLE" // or WAITING
			}
		} else {
			// Fallback to auto-detection
			detected := DetectState(p.TTY)
			desc = detected
			if detected != "IDLE" {
				statusMarker = "BUSY"
			} else {
				statusMarker = "IDLE"
			}
		}

		target := fmt.Sprintf("%s:%s.%s", p.Session, p.WindowIndex, p.PaneIndex)

		if !showIdle && statusMarker == "IDLE" {
			continue
		}

		if desc == "IDLE" {
			desc = "Shell"
		}

		fmt.Printf("[%s] - %s (%s) ||| %s\n",
			branch,
			desc,
			statusMarker,
			target,
		)
	}
}

func getSelfTTY() (string, error) {
	cmd := exec.Command("tty")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
