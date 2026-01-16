package main

import (
	"fmt"
	"strconv"
	"taskmaster-daemon/executor"
)

type Cmd struct {
	base       string
	flags      []string
	profile_id int
	err        bool
}

// Parses the cmd and do all the data cleanup
func (c *Cmd) Parse_cmd(msg *Msg) {
	c.base = msg.get_cmd()
	fmt.Println("base =", c.base)
	c.flags = msg.get_flags()
	c.profile_id = msg.get_profile_id()
	msg.clean_content()
}

// Error sender
func (c *Cmd) send_error(msg *Msg, errorStr string) {
	msg.add_payload("cmd", "error")
	msg.add_payload("flags", errorStr)
}

// Execute commands
func (c *Cmd) Execute(config []File_Config, manager *executor.Manager, msg *Msg) {
	switch c.base {
	case "load":
		if c.flags == nil {
			c.send_error(msg, "Load missing target")
			return
		}
		tmp := get_config_from_file_name(c.flags[0])
		if tmp == nil {
			c.send_error(msg, "Check file existance")
			return
		}
		execConfig := convertToExecutorConfig(*tmp)
		newProfileID := manager.AddProfile(execConfig)
		tasks, err := manager.InfoStatusTasks(newProfileID)
		if err != nil {
			c.send_error(msg, err.Error())
			return
		}

		msg.add_payload("cmd", "load")
		msg.add_payload("flags", c.flags[0])
		msg.add_payload("id", newProfileID)
		msg.add_payload("task", tasks)

	case "reload":
		// Relauch a profile (stop it, reread the config file, launch it again)
		if c.flags == nil {
			c.send_error(msg, "Reload missing target")
			return
		}
		tmp := get_config_from_file_name(c.flags[0])
		if tmp == nil {
			c.send_error(msg, "Check file existance")
			return
		}
		PrintFile_ConfigStruct(*tmp) //temporary
		execConfig := convertToExecutorConfig(*tmp)
		profileID, _ := strconv.Atoi(c.flags[0])
		newProfileID, err := manager.ReloadProfile(execConfig, profileID)
		tasks, err := manager.InfoStatusTasks(newProfileID)
		if err != nil {
			c.send_error(msg, err.Error())
			return
		}

		msg.add_payload("cmd", "reload")
		msg.add_payload("flags", c.flags[0])
		msg.add_payload("id", newProfileID)
		msg.add_payload("task", tasks)

	case "stop":
		if c.flags == nil {
			c.send_error(msg, "Stop missing target")
			return
		}
		if c.profile_id == 0 {
			c.send_error(msg, "Lack connection to profile")
			return
		}
		taskID, _ := strconv.Atoi(c.flags[0])
		//TODO check if is number?
		newProfileID, err := manager.Stop(c.profile_id, taskID)
		if err != nil {
			c.send_error(msg, err.Error())
			return
		}

		msg.add_payload("id", newProfileID)

	case "start":
		if c.flags == nil {
			c.send_error(msg, "Start missing target")
			return
		}
		if c.profile_id == 0 {
			c.send_error(msg, "Lack connection to profile")
			return
		}
		taskID, _ := strconv.Atoi(c.flags[0])
		newProfileID, err := manager.Start(c.profile_id, taskID)
		if err != nil {
			c.send_error(msg, err.Error())
			return
		}

		msg.add_payload("id", newProfileID)

	case "restart":
		if c.flags == nil {
			c.send_error(msg, "Restart missing target")
			return
		}
		if c.profile_id == 0 {
			c.send_error(msg, "Lack connection to profile")
			return
		}
		taskID, _ := strconv.Atoi(c.flags[0])
		newProfileID, err := manager.Restart(c.profile_id, taskID)
		if err != nil {
			c.send_error(msg, err.Error())
			return
		}

		msg.add_payload("id", newProfileID)

	case "kill":
		if c.flags == nil {
			c.send_error(msg, "Kill missing target")
			return
		}
		if c.profile_id == 0 {
			c.send_error(msg, "Lack connection to profile")
			return
		}
		taskID, _ := strconv.Atoi(c.flags[0])
		newProfileID, err := manager.Kill(c.profile_id, taskID)
		if err != nil {
			c.send_error(msg, err.Error())
			return
		}

		msg.add_payload("id", newProfileID)

	case "describe":
		if c.flags == nil {
			c.send_error(msg, "Describe missing target")
			return
		}
		if c.profile_id == 0 {
			c.send_error(msg, "Lack connection to profile")
			return
		}
		taskID, _ := strconv.Atoi(c.flags[0])
		taskDetail, err := manager.DescribeTask(c.profile_id, taskID)
		if err != nil {
			c.send_error(msg, err.Error())
			return
		}

		msg.add_payload("task", taskDetail)

	case "ps":
		// List profiles
		msg.add_payload("profiles", manager.ListProfiles())

	case "ls":
		if c.profile_id == 0 {
			c.send_error(msg, "Lack connection to profile")
			return
		}
		// List all tasks of a profile
		tasks, err := manager.InfoStatusTasks(c.profile_id)
		if err != nil {
			c.send_error(msg, err.Error())
			return
		}

		msg.add_payload("procs", tasks)

	case "ch":
		if c.flags == nil {
			c.send_error(msg, "ch missing targets")
			return
		}
		// Change current profile
		profileID, _ := strconv.Atoi(c.flags[0])
		exists, err := manager.CheckProfileExists(profileID)
		if err != nil {
			c.send_error(msg, err.Error())
			return
		}

		msg.add_payload("id", exists)

	default:
		c.send_error(msg, "Da hell is that? "+c.base)
	}
}

func (c *Cmd) empty_cmd() {
	c.base = ""
	c.flags = nil
}
