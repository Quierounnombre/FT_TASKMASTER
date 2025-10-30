package executor

import (
	"fmt"
	"sync"
)

type Profile struct {
	ID       int
	executor *Executor
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

func (m *Manager) AddProfile() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	profileID := m.nextProfile
	m.nextProfile++

	m.profiles[profileID] = &Profile{
		ID:       profileID,
		executor: NewExecutor(),
	}
	return profileID
}

func (m *Manager) Execute(profileID int, name, logPath, command string, args ...string) (int, error) {
	m.mu.Lock()
	profile, exists := m.profiles[profileID]
	if !exists {
		m.mu.Unlock()
		return 0, fmt.Errorf("profile %d not found", profileID)
	}

	taskID := m.nextID
	m.nextID++
	m.mu.Unlock()

	err := profile.executor.Start(taskID, name, logPath, command, args...)
	if err != nil {
		return 0, err
	}
	return taskID, nil
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

func (m *Manager) GetStatus(profileID, taskID int) (Status, error) {
	m.mu.RLock()
	profile, exists := m.profiles[profileID]
	m.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("profile %d not found", profileID)
	}
	return profile.executor.GetStatus(taskID)
}

func (m *Manager) ListProfiles() []int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]int, 0, len(m.profiles))
	for id := range m.profiles {
		ids = append(ids, id)
	}
	return ids
}
