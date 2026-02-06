package executor

import (
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
}

func (w *Watcher) Stop() {
	close(w.stopChan)
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

func (w *Watcher) checkProfile(profile *Profile) {
	profile.executor.mu.Lock()
	defer profile.executor.mu.Unlock()

	for id, task := range profile.executor.tasks {
		if task.Status == StatusTerminating {
			continue
		}
		if task.Status == StatusPending {
			profile.executor.mu.Unlock() // Do we really need this?
			w.manager.Start(profile.ID, id)
			profile.executor.mu.Lock()
			continue
		}
		if task.Status == StatusRunning && task.Cmd.Process != nil {
			if err := task.Cmd.Process.Signal(syscall.Signal(0)); err != nil {
				w.handleProcessDeath(task, profile.executor)
			}
		} else if task.Status == StatusFailed {
			w.handleRestart(profile.ID, id, task, profile.executor)
		}
	}
}

func (w *Watcher) handleProcessDeath(task *Task, executor *Executor) {
	if task.Cmd.ProcessState != nil {
		if status, ok := task.Cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			task.ExitCode = status.ExitStatus()
		}
	}

	if executor.isExpectedExitCode(task) {
		task.Status = StatusStopped
	} else {
		task.Status = StatusFailed
	}
}

func (w *Watcher) handleRestart(profileID, taskID int, task *Task, executor *Executor) {
	policy := task.restartPolicy
	maxRestarts := task.MaxRestarts

	shouldRestart := false
	switch policy {
	case "always":
		shouldRestart = true
	case "unexpected":
		shouldRestart = !executor.isExpectedExitCode(task)
	}

	if shouldRestart && (maxRestarts == 0 || task.RestartCount < maxRestarts) {
		if task.launchWait > 0 {
			time.Sleep(task.launchWait)
		}
		task.RestartCount++
		task.Status = StatusPending
		executor.mu.Unlock()
		w.manager.Start(profileID, taskID)
		executor.mu.Lock()
	}
}
