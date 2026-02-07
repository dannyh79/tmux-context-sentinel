package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	detectCmd := flag.NewFlagSet("detect", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'detect' or 'list' subcommands")
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
		listCmd.Parse(os.Args[2:])
		runList()
	default:
		fmt.Println("expected 'detect' or 'list' subcommands")
		os.Exit(1)
	}
}

func runDetect(tty string) {
	state := DetectState(tty)
	fmt.Println(state)
}

func runList() {
	panes, err := GetPanes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing panes: %v\n", err)
		os.Exit(1)
	}

	for _, p := range panes {
		branch := GetGitBranch(p.Path)
		state := DetectState(p.TTY)
		
		// Format: [Session:Window:Pane] [Branch] - [Description] (BUSY/IDLE)
		// We want the output to be easily parseable by fzf and useful for switching.
		// For switching, we need the target. fzf can return the whole line, we can parse it in shell.
		
		// Target ID for tmux switch-client: session:window.pane_index
		// Actually switch-client takes session, or -t target-pane. 
		// If we use switch-client -t session:window.pane, it switches.
		
		target := fmt.Sprintf("%s:%s.%s", p.Session, p.WindowIndex, p.PaneIndex)
		
		// Icon/State logic
		statusMarker := "IDLE"
		if state != "IDLE" {
			statusMarker = "BUSY"
		}

		// Description: if BUSY, show tool. If IDLE, show "Shell" or last command? Just "Shell" is fine or empty.
		desc := state
		if state == "IDLE" {
			desc = "Shell"
		}

		// Pad for alignment? fzf handles it well enough.
		// Output format:
		// target | Display String
		
		// Let's print a line that fzf shows. We can use --with-nth to hide the target id if we want, 
		// or just put it at the end/start. 
		// Best practice: `Display Text...   | target_id` (hidden or extracted)
		// Or: `target_id: Display Text`
		
		fmt.Printf("%-20s [%s] - %s (%s) \t%s\n", 
			fmt.Sprintf("[%s:%s:%s]", p.Session, p.WindowIndex, p.PaneIndex),
			branch,
			desc,
			statusMarker,
			target, // Hidden ID column for the script to pick up
		)
	}
}
