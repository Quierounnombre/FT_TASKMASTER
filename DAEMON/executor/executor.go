package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusKilled  Status = "killed"
	StatusRunning Status = "running"
	StatusStopped Status = "stopped"
	StatusFailed  Status = "failed"
	StatusSuccess Status = "success"
)

type Task struct {
	ID                int
	Name              string
	Cmd               *exec.Cmd
	Status            Status
	ExitCode          int
	RestartCount      int
	MaxRestarts       int
	StartTime         time.Time
	StdoutWriter      io.Writer
	StderrWriter      io.Writer
	Env               []string
	WorkingDir        string
	ExpectedExitCodes []int
	Umask             *int
}

type Executor struct {
	mu    sync.RWMutex
	tasks map[int]*Task
}

func NewExecutor() *Executor {
	return &Executor{
		tasks: make(map[int]*Task),
	}
}

func (e *Executor) InitTask(id int, name, command string, stdout, stderr io.Writer, env []string, workingDir string, expectedExitCodes []int, umask *int, args ...string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	cmd := exec.Command(command, args...)

	task := &Task{
		ID:                id,
		Name:              name,
		Cmd:               cmd,
		Status:            StatusPending,
		StdoutWriter:      stdout,
		StderrWriter:      stderr,
		Env:               env,
		WorkingDir:        workingDir,
		ExpectedExitCodes: expectedExitCodes,
		Umask:             umask,
	}
	e.tasks[id] = task

	return nil
}

func (e *Executor) Start(id int) error {
	e.mu.Lock()
	task, exists := e.tasks[id]
	if !exists {
		e.mu.Unlock()
		return fmt.Errorf("task %d not found", id)
	}

	if task.Status != StatusPending {
		e.mu.Unlock()
		return fmt.Errorf("task %d is not pending", id)
	}

	task.Status = StatusRunning
	task.StartTime = time.Now()

	if task.StdoutWriter != nil {
		task.Cmd.Stdout = task.StdoutWriter
	} else {
		task.Cmd.Stdout = io.Discard
	}
	if task.StderrWriter != nil {
		task.Cmd.Stderr = task.StderrWriter
	} else {
		task.Cmd.Stderr = io.Discard
	}
	if task.Env != nil {
		task.Cmd.Env = task.Env
	}
	if task.WorkingDir != "" {
		task.Cmd.Dir = task.WorkingDir
	}
	if task.Umask != nil {
		if task.Cmd.SysProcAttr == nil {
			task.Cmd.SysProcAttr = &syscall.SysProcAttr{}
		}
		task.Cmd.SysProcAttr.Umask = *task.Umask
	}
	e.mu.Unlock()

	err := task.Cmd.Run()

	e.mu.Lock()
	defer e.mu.Unlock()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				task.ExitCode = status.ExitStatus()
			}
		}
	} else {
		task.ExitCode = 0
	}

	if e.isExpectedExitCode(task) {
		task.Status = StatusStopped
	} else {
		task.Status = StatusFailed
	}

	return nil
}

func (e *Executor) isExpectedExitCode(task *Task) bool {
	if task.ExpectedExitCodes == nil {
		return task.ExitCode == 0
	}
	for _, code := range task.ExpectedExitCodes {
		if task.ExitCode == code {
			return true
		}
	}
	return false
}

func (e *Executor) GetStatus(id int) (Status, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	task, exists := e.tasks[id]
	if !exists {
		return "", fmt.Errorf("task %d not found", id)
	}
	return task.Status, nil
}

func (e *Executor) Stop(id int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	task, exists := e.tasks[id]
	if !exists {
		return fmt.Errorf("task %d not found", id)
	}

	if task.Status != StatusRunning {
		return fmt.Errorf("task %d is not running", id)
	}

	if err := task.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to stop task: %w", err)
	}

	return nil
}

func (e *Executor) Kill(id int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	task, exists := e.tasks[id]
	if !exists {
		return fmt.Errorf("task %d not found", id)
	}

	if task.Status != StatusRunning {
		return fmt.Errorf("task %d is not running", id)
	}

	if err := task.Cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill task: %w", err)
	}

	return nil
}
