package executor

import (
	"fmt"
	"io"
	"os/exec"
	"slices"
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

		// Cmd configuration - run through shell to handle scripts and arguments
		cmd := exec.Command("/bin/sh", "-c", process.Cmd)
		// Set environment variables
		cmd.Env = envSlice

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
			CmdStr:            process.Cmd,
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

// Check if task exists
func (e *Executor) CheckTaskExists(id int) (*Task, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	task, exists := e.tasks[id]
	if !exists {
		return task, fmt.Errorf("task %d not found", id)
	}
	return task, nil
}

func (e *Executor) isExpectedExitCode(task *Task) bool {
	// Checks if exit status is in expected exit codes
	if task.ExpectedExitCodes == nil {
		return true
	}
	return slices.Contains(task.ExpectedExitCodes, task.ExitCode)
}

// recreateCmd creates a fresh exec.Cmd from the task's stored configuration.
// Go's exec.Cmd can only be used once, so this is needed for restarts.
func (e *Executor) recreateCmd(task *Task) {
	cmd := exec.Command("/bin/sh", "-c", task.CmdStr)
	cmd.Env = task.Env
	if task.StdoutWriter != nil {
		cmd.Stdout = task.StdoutWriter
	} else {
		cmd.Stdout = io.Discard
	}
	if task.StderrWriter != nil {
		cmd.Stderr = task.StderrWriter
	} else {
		cmd.Stderr = io.Discard
	}
	if task.WorkingDir != "" {
		cmd.Dir = task.WorkingDir
	}
	task.Cmd = cmd
}

func (e *Executor) Start(id int) (int, error) {
	task, err := e.CheckTaskExists(id)
	if err != nil {
		return -1, err
	}

	e.mu.Lock()
	e.recreateCmd(task)
	task.Status = StatusRunning
	task.StartTime = time.Now()
	e.mu.Unlock()

	err = task.Cmd.Run()

	e.mu.Lock()
	defer e.mu.Unlock()

	// Sync and close file handles if they are *os.File
	if f, ok := task.StdoutWriter.(interface{ Sync() error }); ok {
		f.Sync()
	}
	if f, ok := task.StderrWriter.(interface{ Sync() error }); ok {
		f.Sync()
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				task.ExitCode = status.ExitStatus()
			}
		}
		task.Status = StatusFailed
		return -1, err
	} else {
		task.Status = StatusSuccess
	}

	return id, nil
}

func (e *Executor) GetStatus(id int) (Status, error) {
	task, err := e.CheckTaskExists(id)
	if err != nil {
		return "", err
	}
	return task.Status, nil
}

func (e *Executor) GetTaskDetail(id int) (*TaskDetail, error) {
	task, err := e.CheckTaskExists(id)
	if err != nil {
		return nil, err
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
		RestartPolicy:     task.restartPolicy,
	}
	return taskDetail, nil
}

func (e *Executor) Stop(id int) (int, error) {

	task, err := e.CheckTaskExists(id)
	if err != nil {
		return -1, err
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	if task.Status != StatusRunning {
		return -1, fmt.Errorf("task %d is not running", id)
	}

	if err := task.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return -1, fmt.Errorf("failed to stop task: %w", err)
	}

	return id, nil
}

func (e *Executor) Kill(id int) (int, error) {
	task, err := e.CheckTaskExists(id)
	if err != nil {
		return -1, err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if task.Status != StatusRunning {
		return -1, fmt.Errorf("task %d is not running", id)
	}

	if err := task.Cmd.Process.Kill(); err != nil {
		return -1, fmt.Errorf("failed to kill task: %w", err)
	}

	return id, nil
}

// Restart kills the task and sets it to StatusPending.
// The watcher will detect the pending status and start it asynchronously,
// avoiding blocking the daemon's main event loop.
func (e *Executor) Restart(id int) (int, error) {
	task, err := e.CheckTaskExists(id)
	if err != nil {
		return -1, err
	}

	e.mu.Lock()

	if task.Status == StatusRunning && task.Cmd.Process != nil {
		task.Cmd.Process.Kill()
	}
	task.Status = StatusPending

	e.mu.Unlock()

	return id, nil
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

func (e *Executor) InfoStatusTasks() ([]*TaskInfo, error) {
	taskIDs := e.ListTasks()
	tasksInfo := make([]*TaskInfo, 0, len(taskIDs))

	for _, taskID := range taskIDs {
		info, err := e.GetTaskInfo(taskID)
		if err != nil {
			return nil, err
		}
		tasksInfo = append(tasksInfo, info)
	}
	return tasksInfo, nil
}

func (e *Executor) GetTaskInfo(id int) (*TaskInfo, error) {
	task, err := e.CheckTaskExists(id)
	if err != nil {
		return nil, err
	}

	taskInfo := &TaskInfo{
		TaskID:      task.ID,
		Name:        task.Name,
		Status:      task.Status,
		TimeRunning: time.Since(task.StartTime).String(),
	}
	return taskInfo, nil
}
