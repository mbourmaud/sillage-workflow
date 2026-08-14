package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mbourmaud/sillage-workflow/internal/project"
	"github.com/mbourmaud/sillage-workflow/internal/workflow"
)

var version = "0.1.0-dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "doctor":
		return runDoctor(args[1:], stdout, stderr)
	case "transition":
		return runTransition(args[1:], stdout, stderr)
	case "version":
		fmt.Fprintf(stdout, "sillage %s\n", version)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		printUsage(stderr)
		return 2
	}
}

func runDoctor(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("doctor", flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("root", ".", "project root")
	jsonOutput := flags.Bool("json", false, "print JSON")
	if err := flags.Parse(args); err != nil {
		return 2
	}

	report := project.Inspect(*root)
	if *jsonOutput {
		if err := json.NewEncoder(stdout).Encode(report); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
	} else if report.OK {
		fmt.Fprintln(stdout, "project contract: ready")
	} else {
		for _, finding := range report.Findings {
			fmt.Fprintf(stdout, "%s: %s\n", finding.Code, finding.Path)
		}
	}
	if !report.OK {
		return 1
	}
	return 0
}

func runTransition(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("transition", flag.ContinueOnError)
	flags.SetOutput(stderr)
	taskPath := flags.String("task", "", "task JSON path")
	target := flags.String("to", "", "target status")
	jsonOutput := flags.Bool("json", false, "print JSON")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *taskPath == "" || *target == "" {
		fmt.Fprintln(stderr, "transition requires --task and --to")
		return 2
	}

	content, err := os.ReadFile(*taskPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	var task workflow.Task
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&task); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		fmt.Fprintln(stderr, "task must contain exactly one JSON object")
		return 2
	}
	contract := workflow.ValidateTask(task)
	if !contract.OK {
		writeTransitionResult(stdout, stderr, contract, *jsonOutput)
		return 1
	}
	result := workflow.ValidateTransition(task, workflow.Status(*target))
	if !writeTransitionResult(stdout, stderr, result, *jsonOutput) {
		return 2
	}
	if !result.OK {
		return 1
	}
	return 0
}

func writeTransitionResult(stdout io.Writer, stderr io.Writer, result workflow.TransitionResult, jsonOutput bool) bool {
	if jsonOutput {
		if err := json.NewEncoder(stdout).Encode(result); err != nil {
			fmt.Fprintln(stderr, err)
			return false
		}
	} else {
		fmt.Fprintln(stdout, result.Code)
	}
	return true
}

func printUsage(writer io.Writer) {
	fmt.Fprintln(writer, "usage: sillage <doctor|transition|version>")
}
