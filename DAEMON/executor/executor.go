package executor

import (
	"fmt"
	"io"
	"os/exec"
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

func (e *Executor) Start(id int) int {
	e.mu.Lock()
	task, exists := e.tasks[id]
	if !exists {
		e.mu.Unlock()
		fmt.Errorf("task %d not found", id)
		return -1
	}

	task.Status = StatusPending
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
		return -1
	} else {
		task.Status = StatusSuccess
	}

	return id
}

func (e *Executor) isExpectedExitCode(task *Task) bool {
	// Checks if exit status is in expected exit codes
	if task.ExpectedExitCodes == nil {
		return true
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

func (e *Executor) GetTaskDetail(id int) *TaskDetail {
	e.mu.RLock()
	defer e.mu.RUnlock()

	task, exists := e.tasks[id]
	if !exists {
		fmt.Errorf("task %d not found", id)
		return nil
	}

	startTimeStr := ""
	if !task.StartTime.IsZero() {
		startTimeStr = task.StartTime.Format("2006-01-02T15:04:05Z07:00")
	}
	
	taskDetail := &TaskDetail{
		ID:                task.ID,
		Name:              task.Name,
		Cmd:               task.Cmd.String(),
		Status:            task.Status,
		ExitCode:          task.ExitCode,
		RestartCount:      task.RestartCount,
		MaxRestarts:       task.MaxRestarts,
		StartTime:         startTimeStr,
		Env:               task.Env,
		WorkingDir:        task.WorkingDir,
		ExpectedExitCodes: task.ExpectedExitCodes,
		Umask:             task.Umask,
		RestartPolicy:	   task.restartPolicy,
	}
	return taskDetail
}

func (e *Executor) Stop(id int) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	task, exists := e.tasks[id]
	if !exists {
		fmt.Errorf("task %d not found", id)
		return -1
	}

	if task.Status != StatusRunning {
		fmt.Errorf("task %d is not running", id)
		return -1
	}

	if err := task.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
		fmt.Errorf("failed to stop task: %w", err)
		return -1
	}

	return id
}

func (e *Executor) Kill(id int) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	task, exists := e.tasks[id]
	if !exists {
		fmt.Errorf("task %d not found", id)
		return -1
	}

	if task.Status != StatusRunning {
		fmt.Errorf("task %d is not running", id)
		return -1
	}

	if err := task.Cmd.Process.Kill(); err != nil {
		fmt.Errorf("failed to kill task: %w", err)
		return -1
	}

	return id
}

func (e *Executor) Restart(id int) int {
	if e.Kill(id) == -1 {
		return -1
	}
	if e.Start(id) == -1 {
		return -1
	}

	return id
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

func (e *Executor) InfoStatusTasks() []*TaskInfo {
	taskIDs := e.ListTasks()
	tasksInfo := make([]*TaskInfo, 0, len(taskIDs))
	
	for _, taskID := range taskIDs {
		info := e.GetTaskInfo(taskID)
		tasksInfo = append(tasksInfo, info)
	}
	return tasksInfo
}

func (e *Executor) GetTaskInfo(id int) *TaskInfo {
	e.mu.RLock()
	defer e.mu.RUnlock()

	task, exists := e.tasks[id]
	if !exists {
		fmt.Errorf("task %d not found", id)
		return nil
	}

	taskInfo := &TaskInfo{
		TaskID: task.ID,
		Name:   task.Name,
		Status: task.Status,
		TimeRunning: time.Since(task.StartTime).String(),
	}
	return taskInfo
}
