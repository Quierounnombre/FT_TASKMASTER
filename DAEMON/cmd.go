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
	c.flags = msg.get_flags()
	c.profile_id = msg.get_profile_id()
	msg.clean_content()
}

// Error sender
func (c *Cmd) send_error(msg *Msg, errorStr string, logger *executor.Logger) {
	logger.Error(errorStr)
	msg.add_payload("cmd", "error")
	msg.add_payload("flags", errorStr)
}

// Execute commands
func (c *Cmd) Execute(manager *executor.Manager, msg *Msg) {
	manager.Logger().Info(fmt.Sprintf("Request %s | flags: %v | profile_id: %d", c.base, c.flags, c.profile_id))
	switch c.base {
	case "load":
		if c.flags == nil {
			c.send_error(msg, "Load missing target", manager.Logger())
			return
		}
		tmp := get_config_from_file_name(c.flags[0])
		if tmp == nil {
			c.send_error(msg, "Check file existance", manager.Logger())
			return
		}
		execConfig := convertToExecutorConfig(*tmp, manager.Logger())
		newProfileID := manager.AddProfile(execConfig)
		tasks, err := manager.InfoStatusTasks(newProfileID)
		if err != nil {
			c.send_error(msg, err.Error(), manager.Logger())
			return
		}

		msg.add_payload("cmd", "load")
		msg.add_payload("flags", c.flags[0])
		msg.add_payload("id", newProfileID)
		msg.add_payload("task", tasks)

	case "reload":
		// Relauch a profile (stop it, reread the config file, launch it again)
		if c.flags == nil {
			c.send_error(msg, "Reload missing target", manager.Logger())
			return
		}
		profileID, _ := strconv.Atoi(c.flags[0])
		tmp := get_config_from_file_name(manager.GetProfilePath(profileID))
		if tmp == nil {
			c.send_error(msg, "Check file existance", manager.Logger())
			return
		}
		PrintFile_ConfigStruct(*tmp)
		execConfig := convertToExecutorConfig(*tmp, manager.Logger())
		newProfileID, err := manager.ReloadProfile(execConfig, profileID)
		tasks, err := manager.InfoStatusTasks(newProfileID)
		if err != nil {
			c.send_error(msg, err.Error(), manager.Logger())
			return
		}

		msg.add_payload("cmd", "reload")
		msg.add_payload("flags", c.flags[0])
		msg.add_payload("id", newProfileID)
		msg.add_payload("task", tasks)

	case "unload":
		if c.flags == nil {
			c.send_error(msg, "Unload missing target", manager.Logger())
			return
		}
		profileID, _ := strconv.Atoi(c.flags[0])
		err := manager.RemoveProfile(profileID)
		if err != nil {
			c.send_error(msg, err.Error(), manager.Logger())
			return
		}

		msg.add_payload("cmd", "unload")
		msg.add_payload("flags", c.flags[0])

	case "stop":
		if c.flags == nil {
			c.send_error(msg, "Stop missing target", manager.Logger())
			return
		}
		taskID, _ := strconv.Atoi(c.flags[0])

		newProfileID, err := manager.Stop(taskID)
		if err != nil {
			c.send_error(msg, err.Error(), manager.Logger())
			return
		}

		msg.add_payload("cmd", "stop")
		msg.add_payload("id", newProfileID)

	case "start":
		if c.flags == nil {
			c.send_error(msg, "Start missing target", manager.Logger())
			return
		}
		taskID, _ := strconv.Atoi(c.flags[0])
		newProfileID, err := manager.Start(taskID)
		if err != nil {
			c.send_error(msg, err.Error(), manager.Logger())
			return
		}

		msg.add_payload("cmd", "start")
		msg.add_payload("id", newProfileID)

	case "restart":
		if c.flags == nil {
			c.send_error(msg, "Restart missing target", manager.Logger())
			return
		}
		taskID, _ := strconv.Atoi(c.flags[0])
		newProfileID, err := manager.Restart(taskID)
		if err != nil {
			c.send_error(msg, err.Error(), manager.Logger())
			return
		}

		msg.add_payload("cmd", "restart")
		msg.add_payload("id", newProfileID)

	case "kill":
		if c.flags == nil {
			c.send_error(msg, "Kill missing target", manager.Logger())
			return
		}
		taskID, _ := strconv.Atoi(c.flags[0])
		newProfileID, err := manager.Kill(taskID)
		if err != nil {
			c.send_error(msg, err.Error(), manager.Logger())
			return
		}

		msg.add_payload("cmd", "kill")
		msg.add_payload("id", newProfileID)

	case "describe":
		if c.flags == nil {
			c.send_error(msg, "Describe missing target", manager.Logger())
			return
		}
		taskID, _ := strconv.Atoi(c.flags[0])
		taskDetail, err := manager.DescribeTask(taskID)
		if err != nil {
			c.send_error(msg, err.Error(), manager.Logger())
			return
		}

		msg.add_payload("cmd", "describe")
		msg.add_payload("task", taskDetail)

	case "ps":
		// List profiles
		profiles := manager.ListProfiles()
		if profiles == nil {
			c.send_error(msg, "No profiles found", manager.Logger())
			return
		}
		msg.add_payload("cmd", "ps")
		msg.add_payload("profiles", profiles)

	case "ls":
		if c.profile_id == 0 {
			c.send_error(msg, "Set a profile first dude", manager.Logger())
			return
		}
		// List all tasks of a profile
		tasks, err := manager.InfoStatusTasks(c.profile_id)
		if err != nil {
			c.send_error(msg, err.Error(), manager.Logger())
			return
		}

		msg.add_payload("cmd", "ls")
		msg.add_payload("procs", tasks)

	case "ch":
		if c.flags == nil {
			c.send_error(msg, "ch missing targets", manager.Logger())
			return
		}
		// Change current profile
		profileID, _ := strconv.Atoi(c.flags[0])
		_, err := manager.CheckProfileExists(profileID)
		if err != nil {
			c.send_error(msg, err.Error(), manager.Logger())
			return
		}

		msg.add_payload("cmd", "ch")
		msg.add_payload("id", profileID)

	case "russian":
		if c.profile_id == 0 {
			c.send_error(msg, "Set a profile first dude", manager.Logger())
			return
		}
		tasks, err := manager.InfoStatusTasks(c.profile_id)
		if err != nil {
			c.send_error(msg, err.Error(), manager.Logger())
			return
		}
		//Fast random
		for _, task := range tasks {
			newProfileID, err := manager.Kill(task.TaskID)
			if err != nil {
				c.send_error(msg, err.Error(), manager.Logger())
				return
			}
			msg.add_payload("cmd", "russian")
			msg.add_payload("unlucky", newProfileID)
		}

	default:
		c.send_error(msg, "Da hell is that? "+c.base, manager.Logger())
	}
}

func (c *Cmd) empty_cmd() {
	c.base = ""
	c.flags = nil
}
