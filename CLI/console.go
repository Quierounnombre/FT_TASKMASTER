package main

import (
	"os"
	"strconv"
	"net"
	"github.com/chzyer/readline"
	"encoding/json"
)

const history_path = "/run/.history"
const start_shell = 1

//Sets the configuration for the console
func	set_config() *readline.Config {
	var config_rl readline.Config

	config_rl.HistoryFile = history_path
	config_rl.Prompt = "<>< "
	if (len(os.Args) > start_shell) {
		config_rl.Prompt = ""
	}
	config_rl.HistoryLimit = 100
	return (&config_rl)
}

//Sets up rl library and starts the console
func	console_start(sk net.Conn, encoder *json.Encoder) {
	var rl			*readline.Instance
	var config_rl	*readline.Config
	var err			error
	var profile_id	int

	config_rl = set_config()
	rl, err = readline.NewEx(config_rl)
	if (err != nil) {
		os.Exit(1)
	}
	go recive_data(sk, rl, &profile_id)
	console(rl, encoder, &profile_id)
}

//Starts the console
func	console(rl *readline.Instance, encoder *json.Encoder, profile_id *int) {
	var err		error
	var line	string
	var cmd		Cmd
	
	if (len(os.Args) > start_shell) {
		cmd.Cmd = os.Args[1]
		cmd.profile_id = 0
		send_data(encoder, &cmd)
	}
	for (true) {
		line, err = rl.Readline()
		if (err != nil) {
			os.Exit(1)
		}
		if (!local_cmds(profile_id, line, rl)) {
			cmd.Cmd = line
			cmd.profile_id = *profile_id
			send_data(encoder, &cmd)
		}
	}
}

//Check for cmds that dosen't need daemon
//Returns true if found a cmd that dosen't need daemon, false otherwise
func	local_cmds(profile_id *int, line string, rl *readline.Instance) bool {
	if (line == "help") {
		rl.Write([]byte(cmd_help()))
		return (true)
	}
	if (line == "wichid") {
		rl.Write([]byte("Current ID: " + strconv.Itoa(*profile_id) + "\n"))
		return (true)
	}
	return(false)
}

const bold = "\033[1m"
const reset = "\033[0m"

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
