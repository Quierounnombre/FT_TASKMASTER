package executor

import (
	"strconv"
	"syscall"
	"time"
)

type Watcher struct {
	manager  *Manager
	stopChan chan struct{}
}

func NewWatcher(manager *Manager) *Watcher {
	return &Watcher{
		manager:  manager,
		stopChan: make(chan struct{}),
	}
}

func (w *Watcher) Start() {
	go w.watch()
	w.manager.Logger().Info("Watcher started")
}

func (w *Watcher) Stop() {
	close(w.stopChan)
	w.manager.Logger().Info("Watcher stopped")
}

func (w *Watcher) watch() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopChan:
			return
		case <-ticker.C:
			w.checkAllProfiles()
		}
	}
}

func (w *Watcher) checkAllProfiles() {
	w.manager.mu.RLock()
	profiles := make([]*Profile, 0, len(w.manager.profiles))
	for _, profile := range w.manager.profiles {
		profiles = append(profiles, profile)
	}
	w.manager.mu.RUnlock()

	for _, profile := range profiles {
		w.checkProfile(profile)
	}
}

func (w *Watcher) launchTask(executor *Executor, taskID int) {
	go executor.Start(taskID)
}

func (w *Watcher) checkProfile(profile *Profile) {
	profile.executor.mu.Lock()
	defer profile.executor.mu.Unlock()

	for id, task := range profile.executor.tasks {
		if task.Status == StatusTerminating {
			continue
		}
		if task.Status == StatusPending {
			w.launchTask(profile.executor, id)
			w.manager.Logger().Info("Task " + strconv.Itoa(id) + " started")
			continue
		}
		if task.Status == StatusRunning && task.Cmd.Process != nil {
			if err := task.Cmd.Process.Signal(syscall.Signal(0)); err != nil {
				w.handleProcessDeath(task, profile.executor)
			}
		}
		if task.Status == StatusFailed {
			w.handleRestart(id, task, profile.executor)
			w.manager.Logger().Info("Task " + strconv.Itoa(id) + " restarted")
		}
		if task.Status == StatusStopping {
			go w.handleStop(task, profile.executor)
			w.manager.Logger().Info("Task " + strconv.Itoa(id) + " stopping")
		}
	}
}

func (w *Watcher) handleStop(task *Task, executor *Executor) {
	go executor.Stop(task.ID)
	task.StartTime = time.Now()
	// count time to kill
	for {
		if task.Status == StatusStopped {
			w.manager.Logger().Info("Task " + strconv.Itoa(task.ID) + " stopped successfully")
			break
		}
		// kill if time is up
		if time.Since(task.StartTime) > task.Kill_wait {
			executor.Kill(task.ID)
			w.manager.Logger().Info("Task " + strconv.Itoa(task.ID) + " killed due to timeout")
		}
		time.Sleep(1 * time.Second)
	}
}

func (w *Watcher) handleProcessDeath(task *Task, executor *Executor) {
	if task.Cmd.ProcessState != nil {
		if status, ok := task.Cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			task.ExitCode = status.ExitStatus()
		}
	}

	id := task.ID
	if executor.isExpectedExitCode(task) {
		task.Status = StatusSuccess
		w.manager.Logger().Info("Task " + strconv.Itoa(id) + " finished successfully")
	} else {
		task.Status = StatusFailed
		w.manager.Logger().Info("Task " + strconv.Itoa(id) + " finished with error")
	}
}

func (w *Watcher) handleRestart(taskID int, task *Task, executor *Executor) {
	policy := task.restartPolicy
	maxRestarts := task.MaxRestarts

	// With never the task will not be restarted
	// With on_error the task will be restarted if it exited with an error
	// With always the task will be restarted always
	shouldRestart := false
	switch policy {
	case "always":
		shouldRestart = true
	case "on_error":
		shouldRestart = !executor.isExpectedExitCode(task)
	}

	if shouldRestart && (maxRestarts == 0 || task.RestartCount < maxRestarts) {
		if task.launchWait > 0 {
			time.Sleep(task.launchWait)
		}
		task.RestartCount++
		w.launchTask(executor, taskID)
		w.manager.Logger().Info("Task " + strconv.Itoa(taskID) + " restarted " + strconv.Itoa(task.RestartCount) + " times")
	}
}
