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
	StatusRunning Status = "running"
	StatusStopped Status = "stopped"
	StatusFailed  Status = "failed"
)

type Task struct {
	ID      string
	Cmd     *exec.Cmd
	Status  Status
	LogFile *os.File
	ExitCode int
	StartTime time.Time
}

type Executor struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

func New() *Executor {
	return &Executor{
		tasks: make(map[string]*Task),
	}
}

func (e *Executor) Execute(id, logPath, command string, args ...string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.tasks[id]; exists {
		return fmt.Errorf("task %s already exists", id)
	}

	logFile, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = io.MultiWriter(logFile, os.Stdout)
	cmd.Stderr = io.MultiWriter(logFile, os.Stderr)

	task := &Task{
		ID:      id,
		Cmd:     cmd,
		Status:  StatusPending,
		LogFile: logFile,
	}
	e.tasks[id] = task

	go e.run(task)
	return nil
}

func (e *Executor) run(task *Task) {
	e.mu.Lock()
	task.Status = StatusRunning
	task.StartTime = time.Now()
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
		task.Status = StatusFailed
	} else {
		task.ExitCode = 0
		task.Status = StatusStopped
	}

	task.LogFile.Close()
}

func (e *Executor) GetStatus(id string) (Status, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	task, exists := e.tasks[id]
	if !exists {
		return "", fmt.Errorf("task %s not found", id)
	}
	return task.Status, nil
}

func (e *Executor) Stop(id string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	task, exists := e.tasks[id]
	if !exists {
		return fmt.Errorf("task %s not found", id)
	}

	if task.Status != StatusRunning {
		return fmt.Errorf("task %s is not running", id)
	}

	if err := task.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to stop task: %w", err)
	}

	return nil
}

func (e *Executor) Kill(id string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	task, exists := e.tasks[id]
	if !exists {
		return fmt.Errorf("task %s not found", id)
	}

	if task.Status != StatusRunning {
		return fmt.Errorf("task %s is not running", id)
	}

	if err := task.Cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill task: %w", err)
	}

	return nil
}


