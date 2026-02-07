package main

import (
	"os/exec"
	"strings"
)

type Pane struct {
	ID          string
	TTY         string
	Path        string
	Session     string
	WindowIndex string
	PaneIndex   string
	WindowName  string
}

func GetPanes() ([]Pane, error) {
	return getPanesSafe()
}

func getPanesSafe() ([]Pane, error) {
	delim := "|||"
	format := "#{pane_id}" + delim + "#{pane_tty}" + delim + "#{pane_current_path}" + delim + "#{session_name}" + delim + "#{window_index}" + delim + "#{pane_index}" + delim + "#{window_name}"
	cmd := exec.Command("tmux", "list-panes", "-a", "-F", format)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var panes []Pane
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if line == "" { continue }
		parts := strings.Split(line, delim)
		if len(parts) < 7 {
			continue
		}
		panes = append(panes, Pane{
			ID:          parts[0],
			TTY:         parts[1],
			Path:        parts[2],
			Session:     parts[3],
			WindowIndex: parts[4],
			PaneIndex:   parts[5],
			WindowName:  parts[6],
		})
	}
	return panes, nil
}

func GetGitBranch(path string) string {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "-"
	}
	return strings.TrimSpace(string(out))
}

func SwitchToPane(target string) error {
	// target format: session:window.pane or similar. 
	// We usually use -t target
	return exec.Command("tmux", "switch-client", "-t", target).Run()
}
