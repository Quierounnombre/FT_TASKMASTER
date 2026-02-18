package main

import (
	"encoding/json"
	"net"
	"os"
	"strconv"
	"strings"
	"github.com/chzyer/readline"
)

const history_path = "/run/.history"
const start_shell = 1

// Sets the configuration for the console
func set_config() *readline.Config {
	var config_rl readline.Config

	config_rl.HistoryFile = history_path
	config_rl.Prompt = "<>< "
	if len(os.Args) > start_shell {
		config_rl.Prompt = ""
	}
	config_rl.HistoryLimit = 100
	return (&config_rl)
}

// Sets up rl library and starts the console
func console_start(sk net.Conn, encoder *json.Encoder) {
	var rl *readline.Instance
	var config_rl *readline.Config
	var err error
	var profile_id int

	config_rl = set_config()
	rl, err = readline.NewEx(config_rl)
	if err != nil {
		os.Exit(1)
	}
	profile_id = set_profile_id(rl)
	go recive_data(sk, rl, &profile_id)
	console(rl, encoder, &profile_id)
}

// Starts the console
func console(rl *readline.Instance, encoder *json.Encoder, profile_id *int) {
	var err error
	var line string
	var flags []string
	var cmd Cmd

	if len(os.Args) > start_shell {
		if os.Args[1] == "help" {
			rl.Write([]byte(cmd_help()))
		} else {
			cmd.Cmd = os.Args[1]
			cmd.Flags = os.Args[2:]
			cmd.Profile_id = *profile_id
			send_data(encoder, &cmd)
		}
	}
	for true {
		line, err = rl.Readline()
		if err != nil {
			rl.Close()
			os.Exit(1)
		}
		if !local_cmds(profile_id, line, rl) {
			flags = strings.Split(line, " ")
			cmd.Cmd = flags[0]
			cmd.Flags = flags[1:]
			cmd.Profile_id = *profile_id
			send_data(encoder, &cmd)
		}
	}
}

// Check for cmds that dosen't need daemon
// Returns true if found a cmd that dosen't need daemon, false otherwise
func local_cmds(profile_id *int, line string, rl *readline.Instance) bool {
	if line == "help" {
		rl.Write([]byte(cmd_help()))
		return (true)
	}
	if line == "wichid" {
		rl.Write([]byte("Current ID: " + strconv.Itoa(*profile_id) + "\n"))
		return (true)
	}
	return (false)
}

func set_profile_id(rl *readline.Instance) int {
	var profile		int
	var args_len	int
	var err			error
	var arg_line	string

	profile = 0
	args_len = len(os.Args)
	args_len--
	if args_len > start_shell {
		for args_len > 0 {
			if os.Args[args_len] == "-p" {
				args_len++
				arg_line = os.Args[args_len]
				profile, err = strconv.Atoi(arg_line)
				if err != nil {
					rl.Write([]byte("-p invalid profile, make sure is a int, setting default(0)\n"))
					return 0
				}
				return profile
			}
			args_len--
		}
		rl.Write([]byte("-p not found, setting default(0)\n"))
	}
	return 0
}

func extract_args(og_args []string) []string {
	var args	[]string
	var val		string
	var pos		int

	for pos, val = range og_args {
		if val == "-p" {
			pos++
		} else {
			args = append(args, val)
		}
	}
	return args
}

const bold = "\033[1m"
const reset = "\033[0m"

func cmd_help() string {
	var str string

	str = string(str + bold + "load" + reset + "	{PATH}		Load a taskmaster.yaml in the provided path\n")
	str = string(str + bold + "reload" + reset + "	(ID)		Reload the configuration file with the given task id\n")
	str = string(str + bold + "stop" + reset + "	{TARGET}	Stop the target process\n")
	str = string(str + bold + "start" + reset + "	{TARGET}	Start the target process\n")
	str = string(str + bold + "restart" + reset + "	{TARGET}	Restart the target process\n")
	str = string(str + bold + "describe" + reset + "{TARGET}	Describe the target process\n")
	str = string(str + bold + "ps" + reset + "					List all the profiles and show their task id, within a profile id\n")
	str = string(str + bold + "ls" + reset + "		(ID)		List all the process within a profile id\n")
	str = string(str + bold + "ch" + reset + "		{ID}		Modify the working profile id\n")
	str = string(str + bold + "wichid" + reset + "				Print the current profile id in console\n")
	str = string(str + bold + "help" + reset + "				Show this info\n")
	str = string(str + bold + "-p" + reset + "					Set the profile id to the indicated number\n")
	str = string(str + "\n Information in {} is a MUST and can't be skipped\n")
	str = string(str + "Information in () is a OPTIONAL and the id is the current id\n")
	str = string(str + "Current id is the the one modify by loading a configuration, or ch\n")
	return (str)
}
