package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStorage(t *testing.T) {
	// Setup
	tmpDir, err := os.MkdirTemp("", "ctx-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "state.json")
	CustomStatePath = testPath

	// Test 1: Load empty state
	state, err := LoadState()
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}
	if len(state.Panes) != 0 {
		t.Errorf("Expected empty state, got %d panes", len(state.Panes))
	}

	// Test 2: Update state
	err = UpdateAgentState("/dev/pts/1", "running", "TestAgent")
	if err != nil {
		t.Fatalf("UpdateAgentState failed: %v", err)
	}

	// Test 3: Reload and verify
	state, err = LoadState()
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}
	if s, ok := state.Panes["/dev/pts/1"]; !ok {
		t.Error("Expected pane /dev/pts/1 to exist")
	} else {
		if s.Status != "running" {
			t.Errorf("Expected status 'running', got '%s'", s.Status)
		}
		if s.Agent != "TestAgent" {
			t.Errorf("Expected agent 'TestAgent', got '%s'", s.Agent)
		}
	}

	// Test 4: Prune state
	// Assume only pts/2 is active now
	err = PruneState([]string{"/dev/pts/2"})
	if err != nil {
		t.Fatalf("PruneState failed: %v", err)
	}

	state, err = LoadState()
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}
	if _, ok := state.Panes["/dev/pts/1"]; ok {
		t.Error("Expected /dev/pts/1 to be pruned")
	}
}
