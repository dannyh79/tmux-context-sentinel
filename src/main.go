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
		showIdle := listCmd.Bool("show-idle", false, "Show idle panes as well")
		listCmd.Parse(os.Args[2:])
		runList(*showIdle)
	default:
		fmt.Println("expected 'detect' or 'list' subcommands")
		os.Exit(1)
	}
}

func runDetect(tty string) {
	state := DetectState(tty)
	fmt.Println(state)
}

func runList(showIdle bool) {
	panes, err := GetPanes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing panes: %v\n", err)
		os.Exit(1)
	}

	for _, p := range panes {
		branch := GetGitBranch(p.Path)
		state := DetectState(p.TTY)
		
		target := fmt.Sprintf("%s:%s.%s", p.Session, p.WindowIndex, p.PaneIndex)
		
		statusMarker := "IDLE"
		if state != "IDLE" {
			statusMarker = "BUSY"
		}

		if !showIdle && statusMarker == "IDLE" {
			continue
		}

		desc := state
		if state == "IDLE" {
			desc = "Shell"
		}

		// The user wants clean context: [Branch] - [Description] (BUSY/IDLE)
		// We add the target ID at the end separated by a delimiter for fzf parsing.
		fmt.Printf("[%s] - %s (%s) ||| %s\n", 
			branch,
			desc,
			statusMarker,
			target,
		)
	}
}
