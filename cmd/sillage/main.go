package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mbourmaud/sillage-workflow/internal/project"
	"github.com/mbourmaud/sillage-workflow/internal/release"
	"github.com/mbourmaud/sillage-workflow/internal/taskstore"
	"github.com/mbourmaud/sillage-workflow/internal/workflow"
)

var version = "0.2.0-dev"

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
	case "context":
		return runContext(args[1:], stdout, stderr)
	case "status":
		return runStatus(args[1:], stdout, stderr)
	case "changelog":
		return runChangelog(args[1:], stdout, stderr)
	case "digest":
		return runDigest(args[1:], stdout, stderr)
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

func runContext(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("context", flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("root", ".", "project root")
	taskPath := flags.String("task", "", "optional task JSON path")
	jsonOutput := flags.Bool("json", false, "print JSON")
	if err := flags.Parse(args); err != nil {
		return 2
	}

	var task *workflow.Task
	if *taskPath != "" {
		loaded, err := readTask(*taskPath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		task = &loaded
	}
	report := project.Context(*root, task)
	if *jsonOutput {
		if err := json.NewEncoder(stdout).Encode(report); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
	} else {
		if report.Project.Ready {
			fmt.Fprintln(stdout, "project: ready")
		} else {
			fmt.Fprintln(stdout, "project: needs attention")
			for _, finding := range report.Project.Findings {
				fmt.Fprintf(stdout, "%s: %s\n", finding.Code, finding.Path)
			}
		}
		if report.Task != nil {
			fmt.Fprintf(stdout, "task %s: %s\n", report.Task.ID, report.Task.Status)
			printDelegation(stdout, report.Task.Delegation)
			if report.Task.NextStatus != "" {
				fmt.Fprintf(stdout, "next: %s (%s)\n", report.Task.NextStatus, report.Task.NextAction)
			} else {
				fmt.Fprintf(stdout, "next: %s\n", report.Task.NextAction)
			}
		}
	}
	if !report.Project.Ready || (report.Task != nil && !report.Task.Valid) {
		return 1
	}
	return 0
}

func runStatus(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("status", flag.ContinueOnError)
	flags.SetOutput(stderr)
	taskPath := flags.String("task", "", "task JSON path")
	jsonOutput := flags.Bool("json", false, "print JSON")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *taskPath == "" {
		fmt.Fprintln(stderr, "status requires --task")
		return 2
	}
	task, err := readTask(*taskPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	report := project.TaskStatus(task)
	if *jsonOutput {
		if err := json.NewEncoder(stdout).Encode(report); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
	} else {
		fmt.Fprintf(stdout, "%s: %s\n", report.Status, report.NextAction)
		printDelegation(stdout, report.Delegation)
		if report.NextStatus != "" {
			fmt.Fprintf(stdout, "next: %s\n", report.NextStatus)
		}
	}
	if !report.Valid {
		return 1
	}
	return 0
}

func runChangelog(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "changelog requires check or extract")
		return 2
	}
	command := args[0]
	flags := flag.NewFlagSet("changelog "+command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "CHANGELOG.md", "changelog path")
	version := flags.String("version", "", "release version, with or without a leading v")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	switch command {
	case "check":
		if err := release.Check(*path, *version); err != nil {
			fmt.Fprintf(stderr, "changelog: %s\n", err)
			return 1
		}
		if *version == "" {
			fmt.Fprintln(stdout, "changelog: ready")
		} else {
			fmt.Fprintf(stdout, "changelog: %s ready\n", *version)
		}
		return 0
	case "extract":
		if *version == "" {
			fmt.Fprintln(stderr, "changelog extract requires --version")
			return 2
		}
		notes, err := release.Extract(*path, *version)
		if err != nil {
			fmt.Fprintf(stderr, "changelog: %s\n", err)
			return 1
		}
		if _, err := io.WriteString(stdout, notes); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		return 0
	default:
		fmt.Fprintf(stderr, "unknown changelog command %q (use check or extract)\n", command)
		return 2
	}
}

func runTransition(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("transition", flag.ContinueOnError)
	flags.SetOutput(stderr)
	taskPath := flags.String("task", "", "task JSON path")
	target := flags.String("to", "", "target status")
	write := flags.Bool("write", false, "write the accepted transition atomically")
	jsonOutput := flags.Bool("json", false, "print JSON")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *taskPath == "" || *target == "" {
		fmt.Fprintln(stderr, "transition requires --task and --to")
		return 2
	}

	task, original, err := readTaskFile(*taskPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	contract := workflow.ValidateTask(task)
	if !contract.OK {
		writeTransitionResult(stdout, stderr, contract, *jsonOutput)
		return 1
	}
	result := workflow.ValidateTransition(task, workflow.Status(*target))
	if !result.OK {
		if !writeTransitionResult(stdout, stderr, result, *jsonOutput) {
			return 2
		}
		return 1
	}
	if *write {
		if err := taskstore.WriteTransition(*taskPath, task, workflow.Status(*target), original); err != nil {
			writeTransitionResult(stdout, stderr, workflow.TransitionResult{Code: "write_failed"}, *jsonOutput)
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	if !writeTransitionResult(stdout, stderr, result, *jsonOutput) {
		return 2
	}
	return 0
}

func runDigest(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("digest", flag.ContinueOnError)
	flags.SetOutput(stderr)
	taskPath := flags.String("task", "", "task JSON path")
	jsonOutput := flags.Bool("json", false, "print JSON")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *taskPath == "" {
		fmt.Fprintln(stderr, "digest requires --task")
		return 2
	}
	task, err := readTask(*taskPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if contract := workflow.ValidateTask(task); !contract.OK {
		writeTransitionResult(stdout, stderr, contract, *jsonOutput)
		return 1
	}
	digest := workflow.DecisionDigest(task)
	if *jsonOutput {
		if err := json.NewEncoder(stdout).Encode(struct {
			DecisionDigest string `json:"decision_digest"`
		}{DecisionDigest: digest}); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
	} else {
		fmt.Fprintln(stdout, digest)
	}
	return 0
}

func readTask(path string) (workflow.Task, error) {
	task, _, err := readTaskFile(path)
	return task, err
}

func readTaskFile(path string) (workflow.Task, []byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return workflow.Task{}, nil, err
	}
	var task workflow.Task
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&task); err != nil {
		return workflow.Task{}, nil, err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return workflow.Task{}, nil, fmt.Errorf("task must contain exactly one JSON object")
	}
	return task, content, nil
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

func printDelegation(writer io.Writer, request *workflow.DelegationRequest) {
	if request == nil {
		return
	}
	required := "optional"
	if request.Required {
		required = "required"
	}
	fmt.Fprintf(writer, "delegation: %s/%s in %s -> %s (%s)\n", request.Mode, request.Role, request.Isolation, request.Return, required)
}

func printUsage(writer io.Writer) {
	fmt.Fprintln(writer, "usage: sillage <doctor|context|status|changelog|digest|transition|version>")
}
