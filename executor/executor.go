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
	ID        int
	Name      string
	Cmd       *exec.Cmd
	Status    Status
	LogFile   *os.File
	ExitCode  int
	StartTime time.Time
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

func (e *Executor) Start(id int, name, logPath, command string, args ...string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	logFile, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = io.MultiWriter(logFile, os.Stdout)
	cmd.Stderr = io.MultiWriter(logFile, os.Stderr)

	task := &Task{
		ID:      id,
		Name:    name,
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
