package main

import (
	"taskmaster-daemon/executor"
	"fmt"
)

type Cmd struct {
	base  string
	flags []string
	err   bool
}

// Parses the cmd and do all the data cleanup
func (c *Cmd) Parse_cmd(msg *Msg) {
	c.base = msg.get_cmd()
	fmt.Println("base =", c.base)
	c.flags = msg.get_flags()
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
		manager.ReloadProfile(*tmp, 0) //LOOK OUT 0 IS NOT A DINAMIC SELECTED PROFILE
		msg.add_payload("cmd", "reload")
		msg.add_payload("flags", c.flags[0])
		msg.add_payload("id", manager.AddProfile(execConfig))
		msg.add_payload("task", task)
	case "stop":
		msg.add_payload("id", manager.Stop(c.flags[0]/*TMP: this will be a msg variable*/, c.flags[0])) // profileID and taskID)
	case "start":
		msg.add_payload("id", manager.Start(c.flags[0]/*TMP: this will be a msg variable*/, c.flags[0])) // profileID and taskID
	case "restart":
		msg.add_payload("id", manager.Restart(c.flags[0]/*TMP: this will be a msg variable*/, c.flags[0])) // profileID and taskID
	case "kill":
		msg.add_payload("id", manager.Kill(c.flags[0]/*TMP: this will be a msg variable*/, c.flags[0])) // profileID and taskID
	case "describe":
		msg.add_payload("task", manager.DescribeTask(c.flags[0]/*TMP: this will be a msg variable*/, c.flags[0])) // profileID and taskID
	case "ps":
		// List profiles
		return (manager.ListProfiles())
	case "ls":
		// List all tasks of a profile
		result, Err := manager.InfoStatusTasks(c.flags[0])
		if Err != nil {
			return (Err)
		}
		return result
	case "ch":
		// Change current profile

	case "help":
		return (cmd_help())
	}
	return ("Wrong cmd")
}

func (c *Cmd) empty_cmd() {
	c.base = ""
	c.flags = nil
}
