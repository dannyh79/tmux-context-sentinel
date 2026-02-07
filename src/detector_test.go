package main

import "testing"

func TestDetectStateFromProcs(t *testing.T) {
	tests := []struct {
		name     string
		procs    []Process
		expected string
	}{
		{
			name: "Idle Shell",
			procs: []Process{
				{Command: "-zsh"},
				{Command: "tmux a"},
			},
			expected: "IDLE",
		},
		{
			name: "Aider Running",
			procs: []Process{
				{Command: "-zsh"},
				{Command: "python3 -m aider"},
			},
			expected: "Aider",
		},
		{
			name: "Gemini CLI",
			procs: []Process{
				{Command: "node /usr/local/bin/gemini"},
			},
			expected: "Gemini",
		},
		{
			name: "Kiro CLI",
			procs: []Process{
				{Command: "/bin/kiro"},
			},
			expected: "Kiro",
		},
		{
			name: "Cursor CLI",
			procs: []Process{
				{Command: "cursor"},
			},
			expected: "Cursor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectStateFromProcs(tt.procs)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
