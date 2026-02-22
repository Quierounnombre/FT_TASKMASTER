package executor

import (
	"io"
	"os/exec"
	"sync"
	"time"
)

// Configuration types
type Process struct {
	Name              string
	Cmd               string
	Restart           string
	Stop_signal       string
	WorkingDir        string
	Stdout            io.Writer
	Stderr            io.Writer
	Env               map[string]string
	Restart_atempts   int
	ExpectedExitCodes []int
	Launch_wait       time.Duration
	Kill_wait         time.Duration
	Start_at_launch   bool
	Umask             int
	Num_procs         int
}

type File_Config struct {
	Process []Process
	Path    string
}

// Task types
type Status string

const (
	StatusPending     Status = "pending"
	StatusNotLaunched Status = "not_launched"
	StatusKilled      Status = "killed"
	StatusRunning     Status = "running"
	StatusStopped     Status = "stopped"
	StatusFailed      Status = "failed"
	StatusSuccess     Status = "success"
	StatusTerminating Status = "terminating"
)

type Task struct {
	ID                int
	Name              string
	Cmd               *exec.Cmd
	CmdStr            string
	Status            Status
	ExitCode          int
	RestartCount      int
	MaxRestarts       int
	StartTime         time.Time
	StdoutWriter      io.Writer
	StderrWriter      io.Writer
	Env               []string
	WorkingDir        string
	ExpectedExitCodes []int
	Umask             int
	restartPolicy     string
	launchWait        time.Duration
	startAtLaunch     bool
}

// Executor type
type Executor struct {
	mu    sync.RWMutex
	tasks map[int]*Task
}

// Manager types
type Profile struct {
	ID             int
	executor       *Executor
	configFilePath string
}

type Manager struct {
	mu          sync.RWMutex
	profiles    map[int]*Profile
	nextProfile int
	nextID      int
	watcher     *Watcher
	logger      *Logger
}

// Return types
type ListProfiles struct {
	ProfileID int    `json:"profileID"`
	FilePath  string `json:"filePath"`
}

type TaskInfo struct {
	TaskID      int    `json:"taskID"`
	Name        string `json:"name"`
	Status      Status `json:"status"`
	TimeRunning string `json:"timeRunning"`
}

type TaskDetail struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Cmd               string   `json:"cmd"`
	Status            Status   `json:"status"`
	ExitCode          int      `json:"exitCode"`
	RestartCount      int      `json:"restartCount"`
	MaxRestarts       int      `json:"maxRestarts"`
	StartTime         string   `json:"startTime"`
	Env               []string `json:"env"`
	WorkingDir        string   `json:"workingDir"`
	ExpectedExitCodes []int    `json:"expectedExitCodes"`
	Umask             int      `json:"umask"`
	RestartPolicy     string   `json:"restartPolicy"`
}
