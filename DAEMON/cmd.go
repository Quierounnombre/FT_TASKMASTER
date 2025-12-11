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
	case "ps":
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

	str = string(str + bold + "load" + reset + "	{PATH}		Load a taskmaster.yaml in the provided path\n")
	str = string(str + bold + "reload" + reset + "	(ID)		Reload the configuration file with the given id\n")
	str = string(str + bold + "stop" + reset + "	{TARGET}	Stop the target process\n")
	str = string(str + bold + "start" + reset + "	{TARGET}	Start the target process\n")
	str = string(str + bold + "restart" + reset + "	{TARGET}	Restart the target process\n")
	str = string(str + bold + "describe" + reset + "{TARGET}	Describe the target process\n")
	str = string(str + bold + "ps" + reset + "					List all the profiles and show their id\n")
	str = string(str + bold + "ls" + reset + "		(ID)		List all the process within a give id\n")
	str = string(str + bold + "ch" + reset + "		{ID}		Modify the working id\n")
	str = string(str + bold + "wichid" + reset + "				Print the current id in console\n")
	str = string(str + bold + "help" + reset + "				Show this info\n")
	str = string(str + "\n Information in {} is a MUST and can't be skipped\n")
	str = string(str + "Information in () is a OPTIONAL and the id is the current id\n")
	str = string(str + "Current id is the the one modify by loading a configuration, or ch\n")
	return (str)
}
