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
	Command string
}

// DetectState returns the detected tool name or "IDLE"
func DetectState(tty string) string {
	procs, err := getProcessesOnTTY(tty)
	if err != nil {
		return "IDLE"
	}
	return detectStateFromProcs(procs)
}

func detectStateFromProcs(procs []Process) string {
	for _, p := range procs {
		cmd := strings.ToLower(p.Command)
		
		// Signatures
		if strings.Contains(cmd, "aider") {
			return "Aider"
		}
		if strings.Contains(cmd, "gemini") { // gemini-cli
			return "Gemini"
		}
		if strings.Contains(cmd, "kiro") { // kiro-cli
			return "Kiro"
		}
		if strings.Contains(cmd, "cursor") { // cursor-cli
			return "Cursor"
		}
	}
	return "IDLE"
}

func getProcessesOnTTY(tty string) ([]Process, error) {
	// ps -t <tty> -o pid,ppid,command
	cmd := exec.Command("ps", "-t", tty, "-o", "pid,ppid,command")
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
		if len(parts) < 3 {
			continue
		}
		// PID is parts[0], PPID is parts[1], Command is remainder
		// But wait, command can have spaces.
		// parts[0] = PID
		// parts[1] = PPID
		// parts[2:] = Command parts
		
		p := Process{
			PID:     parts[0],
			PPID:    parts[1],
			Command: strings.Join(parts[2:], " "),
		}
		procs = append(procs, p)
	}

	return procs, nil
}
