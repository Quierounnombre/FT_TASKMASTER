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
		msg.add_payload("cmd", "load")
		msg.add_payload("flags", c.flags[0])
		msg.add_payload("id", manager.AddProfile(execConfig))
		msg.add_payload("task", task)
	case "reload":
		// Relauch a profile (stop it, reread the config file, launch it again)
		tmp := get_config_from_file_name(c.flags[0])
		PrintFile_ConfigStruct(*tmp)
		msg.add_payload("cmd", "reload")
		msg.add_payload("flags", c.flags[0])
		msg.add_payload("id", manager.ReloadProfile(*tmp, c.flags[0]))
		msg.add_payload("task", task)
	case "stop":
		msg.add_payload("id", manager.Stop(c.profile_id, c.flags[0])) // profileID and taskID)
	case "start":
		msg.add_payload("id", manager.Start(c.profile_id, c.flags[0])) // profileID and taskID
	case "restart":
		msg.add_payload("id", manager.Restart(c.profile_id, c.flags[0])) // profileID and taskID
	case "kill":
		msg.add_payload("id", manager.Kill(c.profile_id, c.flags[0])) // profileID and taskID
	case "describe":
		msg.add_payload("task", manager.DescribeTask(c.profile_id, c.flags[0])) // profileID and taskID
	case "ps":
		// List profiles
		msg.add_payload("profiles", manager.ListProfiles())
	case "ls":
		// List all tasks of a profile
		msg.add_payload("procs", manager.InfoStatusTasks(c.profile_id))
	case "ch":
		// Change current profile
		msg.add_payload("id", manager.CheckProfileExists(c.flags[0])) // profileID
	case "help":
		return (cmd_help())
	default:
		msg.add_payload("error", "Unknown command: "+c.base)
	}
}

func (c *Cmd) empty_cmd() {
	c.base = ""
	c.flags = nil
}
