package main

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

type Process struct {
	PID     string
	PPID    string
	State   string
	Command string
}

// AgentStatus represents the detected agent and its status
type AgentStatus struct {
	Name   string
	IsBusy bool
}

// DetectState returns the detected tool name or "IDLE"
func DetectState(tty string) string {
	status := DetectAgentStatus(tty)
	return status.Name
}

// DetectAgentStatus returns the agent name and whether it is busy
func DetectAgentStatus(tty string) AgentStatus {
	procs, err := getProcessesOnTTY(tty)
	if err != nil {
		return AgentStatus{Name: "IDLE", IsBusy: false}
	}
	return detectStatusFromProcs(procs)
}

func detectStatusFromProcs(procs []Process) AgentStatus {
	for _, p := range procs {
		cmd := strings.ToLower(p.Command)

		var name string
		// Signatures
		if strings.Contains(cmd, "aider") {
			name = "Aider"
		} else if strings.Contains(cmd, "gemini") { // gemini-cli
			name = "Gemini"
		} else if strings.Contains(cmd, "kiro") { // kiro-cli
			name = "Kiro"
		} else if strings.Contains(cmd, "cursor") { // cursor-cli
			name = "Cursor"
		} else if strings.Contains(cmd, "opencode") {
			name = "OpenCode"
		}

		if name != "" {
			// Check if busy (R state usually implies running/processing)
			// S = Sleeping, R = Running, D = Uninterruptible, Z = Zombie, T = Stopped
			// On macOS, state can include +, etc. We look for 'R'.
			isBusy := strings.Contains(p.State, "R") || strings.Contains(p.State, "D") || strings.Contains(p.State, "U")

			return AgentStatus{Name: name, IsBusy: isBusy}
		}
	}
	return AgentStatus{Name: "IDLE", IsBusy: false}
}

func getProcessesOnTTY(tty string) ([]Process, error) {
	// ps -t <tty> -o pid,ppid,state,command
	cmd := exec.Command("ps", "-t", tty, "-o", "pid,ppid,state,command")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var procs []Process
	scanner := bufio.NewScanner(bytes.NewReader(out))
	// Skip header
	if scanner.Scan() {
		// Header line
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}
		// PID is parts[0], PPID is parts[1], State is parts[2], Command is remainder
		// parts[0] = PID
		// parts[1] = PPID
		// parts[2] = State
		// parts[3:] = Command parts

		p := Process{
			PID:     parts[0],
			PPID:    parts[1],
			State:   parts[2],
			Command: strings.Join(parts[3:], " "),
		}
		procs = append(procs, p)
	}

	return procs, nil
}
