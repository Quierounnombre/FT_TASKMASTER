package main

import (
	"fmt"
	"io"
	"os"
	"taskmaster-daemon/executor"
)

func convertToExecutorConfig(mainConfig File_Config) executor.File_Config {
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
			mode := os.FileMode(0666 & *p.Umask)
			fmt.Println("CHECA ESTO -->>> p.Umask = ", p.Umask, " / mode = ", mode)
			if f, err := os.OpenFile(p.Stdout, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode); err == nil {
				f.Chmod(mode)
				fmt.Println("Stdout file opened: ", f.Name())
				execProcess.Stdout = f
			} else {
				fmt.Println("Error opening stdout: ", err)
				execProcess.Stdout = io.Discard
			}
		} else {
			execProcess.Stdout = io.Discard
		}

		if p.Stderr != "" {
			mode := os.FileMode(0666 & *p.Umask)
			if f, err := os.OpenFile(p.Stderr, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode); err == nil {
				f.Chmod(mode)
				fmt.Println("Stderr file opened: ", f.Name())
				execProcess.Stderr = f
			} else {
				fmt.Println("Error opening stderr: ", err)
				execProcess.Stderr = io.Discard
			}
		} else {
			execProcess.Stderr = io.Discard
		}

		execConfig.Process = append(execConfig.Process, execProcess)
	}

	return execConfig
}
