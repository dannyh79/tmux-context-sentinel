package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type AgentState struct {
	Status    string `json:"status"` // "running" or "waiting"
	Agent     string `json:"agent"`
	Timestamp int64  `json:"timestamp"`
}

type GlobalState struct {
	Panes map[string]AgentState `json:"panes"`
}

var (
	stateMutex sync.Mutex
	CustomStatePath string // For testing
)

func getStatePath() (string, error) {
	if CustomStatePath != "" {
		return CustomStatePath, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".tmux-context-sentinel")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "state.json"), nil
}

func LoadState() (*GlobalState, error) {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	path, err := getStatePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &GlobalState{Panes: make(map[string]AgentState)}, nil
		}
		return nil, err
	}

	var state GlobalState
	if err := json.Unmarshal(data, &state); err != nil {
		return &GlobalState{Panes: make(map[string]AgentState)}, nil
	}
	if state.Panes == nil {
		state.Panes = make(map[string]AgentState)
	}
	return &state, nil
}

func SaveState(state *GlobalState) error {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	path, err := getStatePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func UpdateAgentState(tty string, status string, agent string) error {
	state, err := LoadState()
	if err != nil {
		return err
	}

	// Clean up stale entries (optional, but good for hygiene)
	// For now, we just update the current one
	state.Panes[tty] = AgentState{
		Status:    status,
		Agent:     agent,
		Timestamp: time.Now().Unix(),
	}

	return SaveState(state)
}

func PruneState(activeTTYs []string) error {
    state, err := LoadState()
    if err != nil {
        return err
    }
    
    activeMap := make(map[string]bool)
    for _, t := range activeTTYs {
        activeMap[t] = true
    }
    
    dirty := false
    for tty := range state.Panes {
        if !activeMap[tty] {
            delete(state.Panes, tty)
            dirty = true
        }
    }
    
    if dirty {
        return SaveState(state)
    }
    return nil
}
