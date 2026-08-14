// Package taskstore provides explicit, local task-record writes.
package taskstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mbourmaud/sillage-workflow/internal/workflow"
)

// ErrConcurrentModification means the task changed after the caller read it.
var ErrConcurrentModification = errors.New("task changed since it was read")

// ErrSymlinkTarget means a write target is a symbolic link and was refused.
var ErrSymlinkTarget = errors.New("refusing to write through a symbolic link")

// ErrInvalidTask means the task does not satisfy the public task contract.
var ErrInvalidTask = errors.New("task does not satisfy its contract")

// ErrInvalidTransition means the requested transition is not currently allowed.
var ErrInvalidTransition = errors.New("transition is not allowed")

// WriteTransition validates and atomically writes one explicit lifecycle transition.
// expected must be the exact bytes observed by the caller before it made its decision.
func WriteTransition(path string, task workflow.Task, target workflow.Status, expected []byte) error {
	contract := workflow.ValidateTask(task)
	if !contract.OK {
		return fmt.Errorf("%w: %s", ErrInvalidTask, contract.Code)
	}
	transition := workflow.ValidateTransition(task, target)
	if !transition.OK {
		return fmt.Errorf("%w: %s", ErrInvalidTransition, transition.Code)
	}

	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return ErrSymlinkTarget
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("task target is not a regular file")
	}
	current, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !bytes.Equal(current, expected) {
		return ErrConcurrentModification
	}

	next := task
	next.Status = target
	encoded, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		return err
	}
	encoded = append(encoded, '\n')

	dir := filepath.Dir(path)
	temporary, err := os.CreateTemp(dir, ".sillage-task-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := temporary.Chmod(0o600); err != nil {
		_ = temporary.Close()
		return err
	}
	if _, err := temporary.Write(encoded); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}

	// Re-check both the path kind and bytes immediately before replacement. This
	// closes the normal concurrent-edit window without requiring a provider lock.
	latestInfo, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if latestInfo.Mode()&os.ModeSymlink != 0 {
		return ErrSymlinkTarget
	}
	latest, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !bytes.Equal(latest, expected) {
		return ErrConcurrentModification
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return err
	}
	return syncDirectory(dir)
}

func syncDirectory(path string) error {
	directory, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer directory.Close()
	// Directory fsync is a durability improvement on Unix, but is not
	// supported uniformly by every platform. The rename already completed, so
	// a platform-specific fsync error must not report a failed write.
	_ = directory.Sync()
	return nil
}
