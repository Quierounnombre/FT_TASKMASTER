package executor

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func NewExecutor(config *File_Config, nextID *int) *Executor {
	e := &Executor{
		tasks: make(map[int]*Task),
	}
	
	// Create tasks for each process in config
	for _, process := range config.Process {
		e.initTask(process, nextID)
	}
	
	return e
}

func (e *Executor) initTask(process Process, nextID *int) {
	// Mutex lock to protect tasks map
	defer e.mu.Unlock()
	e.mu.Lock()

	// If not defined or num_procs > 1
	numProcs := process.Num_procs
	if numProcs <= 0 {
		numProcs = 1
	}

	// Convert env map to slice
	var envSlice []string
	for key, value := range process.Env {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", key, value))
	}
	
	for i := 0; i < numProcs; i++ {
		taskID := *nextID
		*nextID++
		
		instanceName := process.Name
		if numProcs > 1 {
			instanceName = fmt.Sprintf("%s_%d", process.Name, i)
		}

		// Cmd configuration
		cmd := exec.Command(process.Cmd, envSlice...)
		// Stdout and Stderr redirection
		if process.Stdout != nil {
			cmd.Stdout = process.Stdout
		} else {
			cmd.Stdout = io.Discard
		}
		if process.Stderr != nil {
			cmd.Stderr = process.Stderr
		} else {
			cmd.Stderr = io.Discard
		}
		// Working Directory
		if process.WorkingDir != "" {
			cmd.Dir = process.WorkingDir
		}

		// Create and store the task
		task := &Task{
			ID:                taskID,
			Name:              instanceName,
			Cmd:               cmd,
			Status:            StatusPending,
			StdoutWriter:      process.Stdout,
			StderrWriter:      process.Stderr,
			Env:               envSlice,
			WorkingDir:        process.WorkingDir,
			ExpectedExitCodes: process.ExpectedExitCodes,
			Umask:             process.Umask,
			MaxRestarts:       process.Restart_atempts,
			restartPolicy:     process.Restart,
			launchWait:        process.Launch_wait,
		}
		e.tasks[taskID] = task
	}
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

func (e *Executor) Restart(id int) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	task, exists := e.tasks[id]

	

	return nil
}

func (e *Executor) ListTasks() []int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	ids := make([]int, 0, len(e.tasks))
	for id := range e.tasks {
		ids = append(ids, id)
	}
	return ids
}

func (e *Executor) GetTaskInfo(id int) (string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	task, exists := e.tasks[id]
	if !exists {
		return "", fmt.Errorf("task %d not found", id)
	}

	cmdStr := strings.Join(task.Cmd.Args, " ")
	taskInfo := TaskInfo{
		TaskID: task.ID,
		Name:   task.Name,
		Status: task.Status,
		Cmd:    cmdStr,
	}

	jsonBytes, err := json.Marshal(taskInfo)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func (e *Executor) GetTaskDetail(id int) (string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	task, exists := e.tasks[id]
	if !exists {
		return "", fmt.Errorf("task %d not found", id)
	}

	cmdStr := strings.Join(task.Cmd.Args, " ")
	startTimeStr := ""
	if !task.StartTime.IsZero() {
		startTimeStr = task.StartTime.Format(time.RFC3339)
	}

	taskDetail := TaskDetail{
		ID:                task.ID,
		Name:              task.Name,
		Cmd:               cmdStr,
		Status:            task.Status,
		ExitCode:          task.ExitCode,
		RestartCount:      task.RestartCount,
		MaxRestarts:       task.MaxRestarts,
		StartTime:         startTimeStr,
		Env:               task.Env,
		WorkingDir:        task.WorkingDir,
		ExpectedExitCodes: task.ExpectedExitCodes,
		Umask:             task.Umask,
	}

	jsonBytes, err := json.Marshal(taskDetail)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
