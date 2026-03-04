package executor

import (
	"fmt"
	"io"
	"os/exec"
	"reflect"
	"slices"
	"strconv"
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
		cmd := exec.Command("/bin/bash", "-c", process.Cmd)
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
		// Initial status
		var initStatus Status
		if process.Start_at_launch {
			if process.Launch_wait > 0 {
				initStatus = StatusWaiting
			} else {
				initStatus = StatusPending
			}
		} else {
			initStatus = StatusNotLaunched
		}
		if process.Launch_wait > 0 {
			initStatus = StatusWaiting
			process.Start_at_launch = true
		}

		// Create and store the task
		task := &Task{
			ID:                taskID,
			Name:              instanceName,
			Cmd:               cmd,
			CmdStr:            process.Cmd,
			Status:            initStatus,
			StdoutWriter:      process.Stdout,
			StderrWriter:      process.Stderr,
			Stop_signal:       process.Stop_signal,
			Env:               envSlice,
			WorkingDir:        process.WorkingDir,
			ExpectedExitCodes: process.ExpectedExitCodes,
			Umask:             process.Umask,
			MaxRestarts:       process.Restart_atempts,
			restartPolicy:     process.Restart,
			launchWait:        process.Launch_wait,
			startAtLaunch:     process.Start_at_launch,
			Kill_wait:         process.Kill_wait,
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
	if task.ExpectedExitCodes == nil || task.Status == StatusKilled {
		return true
	}
	return slices.Contains(task.ExpectedExitCodes, task.ExitCode)
}

func (e *Executor) updateTaskStatus(id int, status Status) {
	e.mu.Lock()
	defer e.mu.Unlock()
	task, err := e.CheckTaskExists(id)
	if err != nil {
		return
	}
	task.Status = status
}

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
	// Put the subprocess in its own process group so signals reach the
	// actual script child, not just the /bin/sh wrapper.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
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

	// If the task was being stopped, the process exit is expected — mark as stopped
	if task.Status == StatusStopping {
		task.Status = StatusStopped
		task.EndTime = time.Now()
		return id, nil // return early — signal-induced exit is not an error
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				task.ExitCode = status.ExitStatus()
			}
		}
		if !e.isExpectedExitCode(task) {
			task.Status = StatusFailed
			task.EndTime = time.Now()
			return -1, err
		}
	}

	task.EndTime = time.Now()

	return id, nil
}

func (e *Executor) GetStatus(id int) (Status, error) {
	task, err := e.CheckTaskExists(id)
	if err != nil {
		return "", err
	}
	if task.Status == StatusRunning && task.launchWait > 0 && time.Since(task.StartTime) >= task.launchWait {
		return StatusSuccess, nil
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

func (e *Executor) Stop(id int, logger *Logger) {
	task, err := e.CheckTaskExists(id)
	if err != nil {
		return
	}

	e.mu.Lock()
	if task.Status != StatusRunning {
		e.mu.Unlock()
		return
	}
	task.Status = StatusStopping

	// Send the stop signal to the entire process group (negative PID) so the
	// signal reaches the actual script, not just the /bin/sh parent process.
	if err := syscall.Kill(-task.Cmd.Process.Pid, syscall.Signal(task.Stop_signal)); err != nil {
		e.mu.Unlock()
		logger.Error("Task {" + strconv.Itoa(task.ID) + "} failed to stop: " + err.Error())
		return
	}
	e.mu.Unlock()

	time.Sleep(500 * time.Millisecond)

	e.mu.Lock()
	if task.Kill_wait == 0 {
		if task.Status != StatusStopped {
			e.mu.Unlock()
			if _, err := e.Kill(task.ID); err != nil {
				logger.Error("Task {" + strconv.Itoa(task.ID) + "} failed to kill: " + err.Error())
				return
			}
			logger.Error("Task {" + strconv.Itoa(task.ID) + "} killed due to stop didn't work")
			return
		}
	} else {
		// Send stop signal and wait for the task to stop
		StartTime := time.Now()
		for {
			if task.Status == StatusStopped {
				break
			}
			// Kill if time is up
			if time.Since(StartTime) > time.Duration(task.Kill_wait) {
				e.mu.Unlock()
				if _, err := e.Kill(task.ID); err != nil {
					logger.Error("Task {" + strconv.Itoa(task.ID) + "} failed to kill: " + err.Error())
					return
				}
				logger.Error("Task {" + strconv.Itoa(task.ID) + "} killed due to timeout")
				return
			}
			e.mu.Unlock()
			time.Sleep(1 * time.Second)
			e.mu.Lock()
			logger.Info("Task {" + strconv.Itoa(task.ID) + "} waiting to be killed " + strconv.Itoa(int(time.Since(StartTime)/time.Second)) + "s")
		}
	}
	e.mu.Unlock()
	logger.Info("Task {" + strconv.Itoa(task.ID) + "} stopped successfully")
}

func (e *Executor) Kill(id int) (int, error) {
	task, err := e.CheckTaskExists(id)
	if err != nil {
		return -1, err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if task.Status != StatusRunning && task.Status != StatusStopping {
		return -1, fmt.Errorf("task %d is not running or stopping", id)
	}

	if task.Cmd.Process == nil {
		return -1, fmt.Errorf("task %d process not started", id)
	}

	// Kill the entire process group so the script child is also killed.
	if err := syscall.Kill(-task.Cmd.Process.Pid, syscall.SIGKILL); err != nil {
		return -1, fmt.Errorf("failed to kill task: %w", err)
	}
	task.Status = StatusKilled
	task.EndTime = time.Now()

	return id, nil
}

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

// GetTaskByName returns tasks matching the given base process name.
// For num_procs > 1, instance names are "name_0", "name_1", etc.
// This returns all tasks whose Name equals exactly `name` or starts with `name_`.
func (e *Executor) GetTasksByBaseName(baseName string) []*Task {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var matches []*Task
	for _, task := range e.tasks {
		if task.Name == baseName || (len(task.Name) > len(baseName)+1 && task.Name[:len(baseName)+1] == baseName+"_") {
			matches = append(matches, task)
		}
	}
	return matches
}

// AddTask adds a single process as new task(s) to the executor.
// This is the public entrypoint used during profile reload for new processes.
func (e *Executor) AddTask(process Process, nextID *int) {
	e.initTask(process, nextID)
}

// RemoveTask stops/kills a running task and removes it from the tasks map.
func (e *Executor) RemoveTask(taskID int, logger *Logger) {
	task, err := e.CheckTaskExists(taskID)
	if err != nil {
		return
	}

	// Stop if running
	if task.Status == StatusRunning || task.Status == StatusStopping {
		if task.Cmd.Process != nil {
			e.Stop(taskID, logger)
		}
	}

	e.mu.Lock()
	delete(e.tasks, taskID)
	e.mu.Unlock()
}

// processChanged checks if the new Process config differs from the existing task.
func processChanged(task *Task, p Process) bool {
	if task.CmdStr != p.Cmd {
		return true
	}
	if task.WorkingDir != p.WorkingDir {
		return true
	}
	if task.restartPolicy != p.Restart {
		return true
	}
	if task.Stop_signal != p.Stop_signal {
		return true
	}
	if task.MaxRestarts != p.Restart_atempts {
		return true
	}
	if task.Kill_wait != p.Kill_wait {
		return true
	}
	if task.launchWait != p.Launch_wait {
		return true
	}
	if task.startAtLaunch != p.Start_at_launch {
		return true
	}
	if task.Umask != p.Umask {
		return true
	}
	if !reflect.DeepEqual(task.ExpectedExitCodes, p.ExpectedExitCodes) {
		return true
	}
	// Compare env: convert process env map to slice and compare
	var newEnv []string
	for key, value := range p.Env {
		newEnv = append(newEnv, fmt.Sprintf("%s=%s", key, value))
	}
	slices.Sort(newEnv)
	existingEnv := make([]string, len(task.Env))
	copy(existingEnv, task.Env)
	slices.Sort(existingEnv)
	if !reflect.DeepEqual(existingEnv, newEnv) {
		return true
	}
	return false
}

// UpdateTask updates an existing task with new Process config.
// If the config changed, it kills/stops the running process and applies the new config.
// If unchanged, it leaves the task (and its running process) untouched.
func (e *Executor) UpdateTask(taskID int, process Process, logger *Logger) {
	task, err := e.CheckTaskExists(taskID)
	if err != nil {
		return
	}

	if !processChanged(task, process) {
		// Nothing changed — leave the task as-is
		return
	}

	logger.Info("Task {" + strconv.Itoa(taskID) + "} config changed, restarting with new config")

	// Stop if currently running
	if task.Status == StatusRunning || task.Status == StatusStopping {
		if task.Cmd.Process != nil {
			e.Stop(taskID, logger)
		}
	}

	// Convert env map to slice
	var envSlice []string
	for key, value := range process.Env {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", key, value))
	}

	// Determine initial status
	var initStatus Status
	if process.Start_at_launch {
		if process.Launch_wait > 0 {
			initStatus = StatusWaiting
		} else {
			initStatus = StatusPending
		}
	} else {
		initStatus = StatusNotLaunched
	}
	if process.Launch_wait > 0 {
		initStatus = StatusWaiting
		process.Start_at_launch = true
	}

	// Update task fields in-place
	e.mu.Lock()
	task.CmdStr = process.Cmd
	task.Env = envSlice
	task.WorkingDir = process.WorkingDir
	task.restartPolicy = process.Restart
	task.Stop_signal = process.Stop_signal
	task.MaxRestarts = process.Restart_atempts
	task.Kill_wait = process.Kill_wait
	task.launchWait = process.Launch_wait
	task.startAtLaunch = process.Start_at_launch
	task.Umask = process.Umask
	task.ExpectedExitCodes = process.ExpectedExitCodes
	task.StdoutWriter = process.Stdout
	task.StderrWriter = process.Stderr
	task.Status = initStatus
	task.ExitCode = 0
	task.RestartCount = 0
	task.StartTime = time.Time{}
	task.EndTime = time.Time{}
	e.mu.Unlock()
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

	var timeRunning string
	taskStatus, _ := e.GetStatus(id)
	switch task.Status {
	case StatusPending, StatusNotLaunched:
		// not started yet — show nothing
		timeRunning = ""
	case StatusRunning:
		// live elapsed since start
		timeRunning = time.Since(task.StartTime).String()
		if task.launchWait > 0 && time.Since(task.StartTime) >= task.launchWait {
			taskStatus = StatusSuccess
		}
	default:
		// finished (success, failed, stopped, killed) — duration at the moment it ended
		if !task.StartTime.IsZero() && !task.EndTime.IsZero() {
			timeRunning = task.EndTime.Sub(task.StartTime).String()
		}
	}

	taskInfo := &TaskInfo{
		TaskID:      task.ID,
		Name:        task.Name,
		Status:      taskStatus,
		TimeRunning: timeRunning,
	}
	return taskInfo, nil
}
