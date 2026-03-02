package executor

import (
	"fmt"
	"io"
	"os/exec"
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
	fmt.Println("Task " + strconv.Itoa(id) + " finished and now here comes the stop")
	if task.Status == StatusStopping {
		fmt.Println("Task " + strconv.Itoa(id) + " stopped successfully")
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
	fmt.Println("Task " + strconv.Itoa(task.ID) + " has status: " + string(task.Status))
	if task.Kill_wait == 0 {
		fmt.Println("--> Yes")
		if task.Status != StatusStopped {
			fmt.Println("--> YEEEES")
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
		fmt.Println("Task " + strconv.Itoa(task.ID) + " debe ser matado con kill_wait")
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
			fmt.Println("Task " + strconv.Itoa(task.ID) + " waiting to be killed " + strconv.Itoa(int(time.Since(StartTime)/time.Second)) + "s")
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
	switch task.Status {
	case StatusPending, StatusNotLaunched:
		// not started yet — show nothing
		timeRunning = ""
	case StatusRunning:
		// live elapsed since start
		timeRunning = time.Since(task.StartTime).String()
	default:
		// finished (success, failed, stopped, killed) — duration at the moment it ended
		if !task.StartTime.IsZero() && !task.EndTime.IsZero() {
			timeRunning = task.EndTime.Sub(task.StartTime).String()
		}
	}

	taskInfo := &TaskInfo{
		TaskID:      task.ID,
		Name:        task.Name,
		Status:      task.Status,
		TimeRunning: timeRunning,
	}
	return taskInfo, nil
}
