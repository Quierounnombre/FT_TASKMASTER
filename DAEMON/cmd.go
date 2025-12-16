package main

import (
	"taskmaster-daemon/executor"
	"fmt"
)

type Cmd struct {
	base		string
	flags		[]string
	profile_id	int
	err			bool
}

// Parses the cmd and do all the data cleanup
func (c *Cmd) Parse_cmd(msg *Msg) {
	c.base = msg.get_cmd()
	fmt.Println("base =", c.base)
	c.flags = msg.get_flags()
	c.profile_id = msg.get_profile_id()
	msg.clean_content()
}

// EXECUTE COMANDS
func (c *Cmd) Execute(config []File_Config, manager *executor.Manager, msg *Msg) {
	switch c.base {
	case "load":
		tmp := get_config_from_file_name(c.flags[0])
		execConfig := convertToExecutorConfig(*tmp)
		tmp_id := manager.AddProfile(execConfig)
		msg.add_payload("cmd", "load")
		msg.add_payload("flags", c.flags[0])
		msg.add_payload("id", tmp_id)
		msg.add_payload("task", task)
	case "reload":
		for index, element := range config {
			config[index] = *get_config_from_file_name(element.Path)
		}
		return ("Configurations reloaded")
	case "stop":
		//CHECK FOR AVAILABLE PROCESS HERE?
		// We would do manager.Stop(profileID, taskID) but this is TMP
		return (string("Stoped " + c.flags[0]))
	case "start":
		//CHECK FOR AVAILABLE PROCESS HERE?
		// We would do manager.Start(profileID, taskID) but this is TMP
		return (string("Started " + c.flags[0]))
	case "restart":
		if len(c.flags) > 0 {
			return (string("Restarted " + c.flags[0]))
		}
		return (string("Restarting all programs"))
	case "describe":
		result, _ := manager.DescribeTask(0, 0)
		return result
	case "ps":
		// List profiles
		return (manager.ListProfiles())
	case "ls":
		result, _ := manager.InfoStatusTasks(0)
		return result
	case "help":
		return (cmd_help())
	}
	return ("Wrong cmd")
}

func (c *Cmd) empty_cmd() {
	c.base = ""
	c.flags = nil
}
