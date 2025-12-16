package main

import (
	"taskmaster-daemon/executor"
	"fmt"
)

const bold = "\033[1m"
const reset = "\033[0m"

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
func (c *Cmd) Execute(config []File_Config, manager *executor.Manager) string {
	switch c.base {
	case "load":
		tmp := get_config_from_file_name(c.flags[0])
		PrintFile_ConfigStruct(*tmp)
		execConfig := convertToExecutorConfig(*tmp)
		manager.AddProfile(execConfig)
		return (string("Loaded " + c.flags[0]))
	case "reload":
		// Relauch a profile (stop it, reread the config file, launch it again)
		tmp := get_config_from_file_name(c.flags[0]) //THERE SHOULD BE A WAY TO TAKE THE FILE PATH
		PrintFile_ConfigStruct(*tmp)
		manager.ReloadProfile(*tmp, 0) //LOOK OUT 0 IS NOT A DINAMIC SELECTED PROFILE
		return ("Configuration reloaded")
	case "stop":
		manager.Stop(profileID, taskID) but this is TMP
		return (string("Stoped " + c.flags[0]))
	case "start":
		manager.Start(profileID, taskID)
		return (string("Started " + c.flags[0]))
	case "restart":
		manager.Restart(profileID, taskID)
		return (string("Restarting all programs"))
	case "describe":
		result, Err := manager.DescribeTask(0, 0)
		if Err != nil {
			return (Err)
		}
		return result
	case "lsPf":
		// List profiles
		return (manager.ListProfiles())
	case "ls":
		// List all tasks of a profile
		result, Err := manager.InfoStatusTasks(0)
		if Err != nil {
			return (Err)
		}
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

// Returns a list with all the current working cmds
func cmd_help() string {
	var str string

	str = string(str + bold + "load" + reset + "	{PATH} load the taskmaster.yaml file in the current dir\n")
	str = string(str + bold + "reload" + reset + "	reload all the configuration files executing in the taskmaster\n")
	str = string(str + bold + "stop" + reset + "	{TARGET} stop the target process\n")
	str = string(str + bold + "start" + reset + "	{TARGET} start the target process\n")
	str = string(str + bold + "restart" + reset + "	(TARGET) reset the target process, or all of them in case there are no target\n")
	str = string(str + bold + "help" + reset + "	Show this info")
	return (str)
}
