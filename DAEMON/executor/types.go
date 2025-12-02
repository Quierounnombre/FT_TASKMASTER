package executor

import (
	"io"
	"time"
)

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
