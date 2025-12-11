package executor

import (
	"encoding/json"
	"fmt"
)

func NewManager() *Manager {
	m := &Manager{
		profiles:    make(map[int]*Profile),
		nextProfile: 1,
		nextID:      1,
	}
	m.watcher = NewWatcher(m)
	m.watcher.Start()
	return m
}

func (m *Manager) Shutdown() {
	if m.watcher != nil {
		m.watcher.Stop()
	}
}

// Profile Management
func (m *Manager) AddProfile(config File_Config) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	profileID := m.nextProfile
	m.nextProfile++

	executor := NewExecutor(&config, &m.nextID)
	m.profiles[profileID] = &Profile{
		ID:       profileID,
		executor: executor,
		configFilePath: config.Path,
	}
	return profileID
}

func (m *Manager) RemoveProfile(profileID int) error {
	m.mu.RLock()
	profile, exists := m.profiles[profileID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("profile %d not found", profileID)
	}

	profile.executor.mu.Lock()
	for _, task := range profile.executor.tasks {
		task.Status = StatusTerminating
	}
	profile.executor.mu.Unlock()

	taskIDs := profile.executor.ListTasks()
	for _, taskID := range taskIDs {
		profile.executor.Stop(taskID)
	}

	m.mu.Lock()
	delete(m.profiles, profileID)
	m.mu.Unlock()

	return nil
}

func (m *Manager) ListProfiles() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	profileIDs := make([]int, 0, len(m.profiles))
	for id := range m.profiles {
		profileIDs = append(profileIDs, id)
	}

	jsonBytes, _ := json.Marshal(profileIDs)
	return string(jsonBytes)
}

// Task Management
func (m *Manager) ListTasks(profileID int) ([]int, error) {
	m.mu.RLock()
	profile, exists := m.profiles[profileID]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("profile %d not found", profileID)
	}
	return profile.executor.ListTasks(), nil
}

func (m *Manager) InfoStatusTasks(profileID int) (string, error) {
	m.mu.RLock()
	profile, exists := m.profiles[profileID]
	m.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("profile %d not found", profileID)
	}

	taskIDs := profile.executor.ListTasks()
	taskInfos := make([]TaskInfo, 0, len(taskIDs))
	
	for _, taskID := range taskIDs {
		profile.executor.mu.RLock()
		task, exists := profile.executor.tasks[taskID]
		if !exists {
			profile.executor.mu.RUnlock()
			continue
		}
		
		taskInfo := TaskInfo{
			TaskID: task.ID,
			Name:   task.Name,
			Status: task.Status,
			Cmd:    task.Cmd.String(),
		}
		taskInfos = append(taskInfos, taskInfo)
		profile.executor.mu.RUnlock()
	}

	jsonBytes, err := json.Marshal(taskInfos)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func (m *Manager) DescribeTask(profileID, taskID int) (string, error) {
	m.mu.RLock()
	profile, exists := m.profiles[profileID]
	m.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("profile %d not found", profileID)
	}
	return profile.executor.GetTaskDetail(taskID)
}

func (m *Manager) GetStatus(profileID, taskID int) (Status, error) {
	m.mu.RLock()
	profile, exists := m.profiles[profileID]
	m.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("profile %d not found", profileID)
	}
	return profile.executor.GetStatus(taskID)
}

func (m *Manager) Start(profileID, taskID int) error {
	m.mu.RLock()
	profile, exists := m.profiles[profileID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("profile %d not found", profileID)
	}
	return profile.executor.Start(taskID)
}

func (m *Manager) Stop(profileID, taskID int) error {
	m.mu.RLock()
	profile, exists := m.profiles[profileID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("profile %d not found", profileID)
	}
	return profile.executor.Stop(taskID)
}

func (m *Manager) Kill(profileID, taskID int) error {
	m.mu.RLock()
	profile, exists := m.profiles[profileID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("profile %d not found", profileID)
	}
	return profile.executor.Kill(taskID)
}
