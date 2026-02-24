package main

import (
	"fmt"
	"io"
	"os"
	"taskmaster-daemon/executor"
)

func convertToExecutorConfig(mainConfig File_Config, logger *executor.Logger) executor.File_Config {
	execConfig := executor.File_Config{
		Path: mainConfig.Path,
	}

	for _, p := range mainConfig.Process {
		execProcess := executor.Process{
			Name:              p.Name,
			Cmd:               p.Cmd,
			Restart:           p.Restart,
			Stop_signal:       p.Stop_signal,
			WorkingDir:        p.Work_dir,
			Env:               p.Env,
			Restart_atempts:   p.Restart_atempts,
			ExpectedExitCodes: p.Expected_exit,
			Launch_wait:       p.Launch_wait,
			Kill_wait:         p.Kill_wait,
			Start_at_launch:   p.Start_at_launch,
			Umask:             *p.Umask,
			Num_procs:         p.Num_procs,
		}

		if p.Stdout != "" {
			mode := os.FileMode(0666 & ^*p.Umask)
			if f, err := os.OpenFile(p.Stdout, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode); err == nil {
				f.Chmod(mode)
				fmt.Println("> Stdout file opened for ", p.Name, " : ", p.Stdout)
				logger.Info("[" + p.Name + "] stdout → file: " + f.Name() + fmt.Sprintf(" (mode %04o)", mode))
				execProcess.Stdout = f
			} else {
				execProcess.Stdout = io.Discard
			}
		} else {
			execProcess.Stdout = io.Discard
		}

		if p.Stderr != "" {
			mode := os.FileMode(0666 & *p.Umask)
			if f, err := os.OpenFile(p.Stderr, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode); err == nil {
				f.Chmod(mode)
				fmt.Println("> Stderr file opened for ", p.Name, " : ", p.Stderr)
				logger.Info("[" + p.Name + "] stderr → file: " + f.Name() + fmt.Sprintf(" (mode %04o)", mode))
				execProcess.Stderr = f
			} else {
				execProcess.Stderr = io.Discard
			}
		} else {
			execProcess.Stderr = io.Discard
		}

		execConfig.Process = append(execConfig.Process, execProcess)
	}

	return execConfig
}
