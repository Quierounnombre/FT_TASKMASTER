package executor

import (
	"fmt"
	"strconv"
)

func NewManager() *Manager {
	logger, err := New("/tmp/taskmaster.log") //<--- Thi shi temporary
	if err != nil {
		panic("Failed to initialize logger file: " + err.Error())
	}
	m := &Manager{
		profiles:    make(map[int]*Profile),
		nextProfile: 1,
		nextID:      1,
		logger:      logger,
	}
	m.watcher = NewWatcher(m)
	m.watcher.Start()
	return m
}

func (m *Manager) Shutdown() {
	m.logger.Info("Shutting down manager")
	if m.watcher != nil {
		m.watcher.Stop()
	}
}

// Check if profile exists
func (m *Manager) CheckProfileExists(profileID int) (*Profile, error) {
	m.mu.RLock()
	profile, exists := m.profiles[profileID]
	m.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("profile %d not found", profileID)
	}
	return profile, nil
}

// Profile Management
func (m *Manager) AddProfile(config File_Config) int {
	m.logger.Info("Adding profile from: " + config.Path)
	m.mu.Lock()
	defer m.mu.Unlock()

	profileID := m.nextProfile
	m.nextProfile++

	executor := NewExecutor(&config, &m.nextID)
	m.profiles[profileID] = &Profile{
		ID:             profileID,
		executor:       executor,
		configFilePath: config.Path,
	}
	return profileID
}

func (m *Manager) RemoveProfile(profileID int) error {
	m.logger.Info("Removing profile " + strconv.Itoa(profileID))
	profile, err := m.CheckProfileExists(profileID)
	if err != nil {
		m.logger.Error("Cmd: RemoveProfile: Profile " + strconv.Itoa(profileID) + " not found")
		return err
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

func (m *Manager) ReloadProfile(config File_Config, profileID int) (int, error) {
	m.logger.Info("Reloading profile " + strconv.Itoa(profileID))
	profile, err := m.CheckProfileExists(profileID)
	if err != nil {
		m.logger.Error("Cmd: ReloadProfile: Profile " + strconv.Itoa(profileID) + " not found")
		return -1, err
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

	newExecutor := NewExecutor(&config, &m.nextID)

	m.mu.Lock()
	profile.executor = newExecutor
	m.mu.Unlock()

	return profileID, nil
}

func (m *Manager) ListProfiles() []ListProfiles {
	m.logger.Info("Listing profiles")
	m.mu.RLock()
	defer m.mu.RUnlock()

	profileIDs := make([]ListProfiles, 0, len(m.profiles))
	for _, profile := range m.profiles {
		profileIDs = append(profileIDs, ListProfiles{
			ProfileID: profile.ID,
			FilePath:  profile.configFilePath,
		})
	}
	return profileIDs
}

// Task Management
func (m *Manager) ListTasks(profileID int) ([]int, error) {
	m.logger.Info("Listing tasks of profile " + strconv.Itoa(profileID))
	profile, err := m.CheckProfileExists(profileID)
	if err != nil {
		m.logger.Error("Cmd: ListTasks: Profile " + strconv.Itoa(profileID) + " not found")
		return nil, err
	}

	return profile.executor.ListTasks(), nil
}

func (m *Manager) InfoStatusTasks(profileID int) ([]*TaskInfo, error) {
	m.logger.Info("Listing status of tasks of profile " + strconv.Itoa(profileID))
	profile, err := m.CheckProfileExists(profileID)
	if err != nil {
		m.logger.Error("Cmd: InfoStatusTasks: Profile " + strconv.Itoa(profileID) + " not found")
		return nil, err
	}

	backInfoStatusTasks, err := profile.executor.InfoStatusTasks()
	if err != nil {
		m.logger.Error("Cmd: InfoStatusTasks: " + err.Error())
		return nil, err
	}
	return backInfoStatusTasks, nil
}

func (m *Manager) DescribeTask(profileID, taskID int) (*TaskDetail, error) {
	m.logger.Info("Listing status of task " + strconv.Itoa(taskID) + " of profile " + strconv.Itoa(profileID))
	profile, err := m.CheckProfileExists(profileID)
	if err != nil {
		m.logger.Error("Cmd: DescribeTask: Profile " + strconv.Itoa(profileID) + " not found")
		return nil, err
	}

	backDetail, err := profile.executor.GetTaskDetail(taskID)
	if err != nil {
		m.logger.Error("Cmd: DescribeTask: " + err.Error())
		return nil, err
	}
	return backDetail, nil
}

func (m *Manager) GetStatus(profileID, taskID int) (Status, error) {
	m.logger.Info("Listing status of task " + strconv.Itoa(taskID) + " of profile " + strconv.Itoa(profileID))
	profile, err := m.CheckProfileExists(profileID)
	if err != nil {
		m.logger.Error("Cmd: GetStatus: Profile " + strconv.Itoa(profileID) + " not found")
		return "", err
	}

	backStatus, err := profile.executor.GetStatus(taskID)
	if err != nil {
		m.logger.Error("Cmd: GetStatus: " + err.Error())
		return "", err
	}
	return backStatus, nil
}

func (m *Manager) Start(profileID, taskID int) (int, error) {
	m.logger.Info("Starting task " + strconv.Itoa(taskID) + " of profile " + strconv.Itoa(profileID))
	profile, err := m.CheckProfileExists(profileID)
	if err != nil {
		m.logger.Error("Cmd: Start: Profile " + strconv.Itoa(profileID) + " not found")
		return -1, err
	}

	backInt, err := profile.executor.Start(taskID)
	if err != nil {
		m.logger.Error("Cmd: Start: " + err.Error())
		return -1, err
	}
	return backInt, nil
}

func (m *Manager) Stop(profileID, taskID int) (int, error) {
	m.logger.Info("Stopping task " + strconv.Itoa(taskID) + " of profile " + strconv.Itoa(profileID))
	profile, err := m.CheckProfileExists(profileID)
	if err != nil {
		m.logger.Error("Cmd: Stop: Profile " + strconv.Itoa(profileID) + " not found")
		return -1, err
	}

	backInt, err := profile.executor.Stop(taskID)
	if err != nil {
		m.logger.Error("Cmd: Stop: " + err.Error())
		return -1, err
	}
	return backInt, nil
}

func (m *Manager) Kill(profileID, taskID int) (int, error) {
	m.logger.Info("Killing task " + strconv.Itoa(taskID) + " of profile " + strconv.Itoa(profileID))
	profile, err := m.CheckProfileExists(profileID)
	if err != nil {
		m.logger.Error("Cmd: Kill: Profile " + strconv.Itoa(profileID) + " not found")
		return -1, err
	}

	backInt, err := profile.executor.Kill(taskID)
	if err != nil {
		m.logger.Error("Cmd: Kill: " + err.Error())
		return -1, err
	}
	return backInt, nil
}

func (m *Manager) Restart(profileID, taskID int) (int, error) {
	m.logger.Info("Restarting task " + strconv.Itoa(taskID) + " of profile " + strconv.Itoa(profileID))
	profile, err := m.CheckProfileExists(profileID)
	if err != nil {
		m.logger.Error("Cmd: Restart: Profile " + strconv.Itoa(profileID) + " not found")
		return -1, err
	}

	backInt, err := profile.executor.Restart(taskID)
	if err != nil {
		m.logger.Error("Cmd: Restart: " + err.Error())
		return -1, err
	}
	return backInt, nil
}
