package executor

import (
	"fmt"
	"strings"
	"sync"
)

type Profile struct {
	ID       int
	executor *Executor
	configFilePath string
}

type Manager struct {
	mu          sync.RWMutex
	profiles    map[int]*Profile
	nextProfile int
	nextID      int
}

func NewManager() *Manager {
	return &Manager{
		profiles:    make(map[int]*Profile),
		nextProfile: 1,
		nextID:      1,
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
	// TODO: Implement profile removal logic
	// Have all tasks stopped before removing
	return nil
}

func (m *Manager) ListProfiles() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.profiles) == 0 {
		return "[]"
	}

	result := "["
	first := true
	for id := range m.profiles {
		if !first {
			result += ","
		}
		result += fmt.Sprintf("%d", id)
		first = false
	}
	result += "]"
	return result
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
	if len(taskIDs) == 0 {
		return "[]", nil
	}

	jsonStrings := make([]string, 0, len(taskIDs))
	for _, taskID := range taskIDs {
		info, err := profile.executor.GetTaskInfo(taskID)
		if err != nil {
			return "", err
		}
		jsonStrings = append(jsonStrings, info)
	}

	return "[" + strings.Join(jsonStrings, ",") + "]", nil
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
