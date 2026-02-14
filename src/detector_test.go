package main

import "testing"

func TestDetectStatusFromProcs(t *testing.T) {
	tests := []struct {
		name     string
		procs    []Process
		expected AgentStatus
	}{
		{
			name: "Idle Shell",
			procs: []Process{
				{Command: "-zsh", State: "S"},
				{Command: "tmux a", State: "S"},
			},
			expected: AgentStatus{Name: "IDLE", IsBusy: false},
		},
		{
			name: "Aider Running (Idle)",
			procs: []Process{
				{Command: "-zsh", State: "S"},
				{Command: "python3 -m aider", State: "S"},
			},
			expected: AgentStatus{Name: "Aider", IsBusy: false},
		},
		{
			name: "Aider Running (Busy)",
			procs: []Process{
				{Command: "-zsh", State: "S"},
				{Command: "python3 -m aider", State: "R"},
			},
			expected: AgentStatus{Name: "Aider", IsBusy: true},
		},
		{
			name: "Gemini CLI",
			procs: []Process{
				{Command: "node /usr/local/bin/gemini", State: "S"},
			},
			expected: AgentStatus{Name: "Gemini", IsBusy: false},
		},
		{
			name: "Kiro CLI (Busy)",
			procs: []Process{
				{Command: "/bin/kiro", State: "R+"},
			},
			expected: AgentStatus{Name: "Kiro", IsBusy: true},
		},
		{
			name: "Cursor CLI",
			procs: []Process{
				{Command: "cursor", State: "S"},
			},
			expected: AgentStatus{Name: "Cursor", IsBusy: false},
		},
		{
			name: "OpenCode CLI",
			procs: []Process{
				{Command: "opencode run", State: "R"},
			},
			expected: AgentStatus{Name: "OpenCode", IsBusy: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectStatusFromProcs(tt.procs)
			if got != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, got)
			}
		})
	}
}
