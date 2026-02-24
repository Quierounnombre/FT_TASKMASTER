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
	var rl				*readline.Instance
	var config_rl		*readline.Config
	var err				error
	var profile_id		int

	config_rl = set_config()
	rl, err = readline.NewEx(config_rl)
	if err != nil {
		os.Exit(1)
	}
	profile_id = set_profile_id(rl)
	console(rl, encoder, &profile_id, sk)
}

// Starts the console
func console(rl *readline.Instance, encoder *json.Encoder, profile_id *int, sk net.Conn) {
	var err error
	var line string
	var flags []string
	var cmd Cmd

	if len(os.Args) > start_shell {
		if os.Args[1] == "help" {
			rl.Write([]byte(cmd_help()))
		} else {
			cmd.Cmd = os.Args[1]
			cmd.Flags = extract_args(os.Args[2:])
			cmd.Profile_id = obtain_profile_id_from_flags(os.Args[2:], *profile_id, rl)
			send_data(encoder, &cmd)
			recive_data(sk, rl, profile_id)
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
			cmd.Flags = extract_args(flags[1:])
			cmd.Profile_id = obtain_profile_id_from_flags(flags[1:], *profile_id, rl)
			send_data(encoder, &cmd)
			recive_data(sk, rl, profile_id)
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
	var err			error
	var val			string
	var extract		bool

	profile = 0
	extract = false
	if len(os.Args) > start_shell {
		for _, val = range os.Args {
			if val == "-p" {
				extract = true
			} else if extract {
				profile, err = strconv.Atoi(val)
				if err != nil {
					rl.Write([]byte("-p invalid profile, make sure is a int, setting default(0)\n"))
					return 0
				}
				return profile
			}
		}
		if extract {
			rl.Write([]byte("-p value not found, setting default(0)\n"))
		}
	}
	return 0
}

//In case of error it set to the current profile id
func obtain_profile_id_from_flags(args []string, profile_id int, rl *readline.Instance) int {
	var val			string
	var extract		bool
	var new_id		int
	var err			error

	extract = false
	for _, val = range args {
		if val == "-p" {
			extract = true
		} else if extract == true{
			new_id, err = strconv.Atoi(val)
			if err != nil {
				rl.Write([]byte("-p invalid profile, make sure is a int, setting current(" + strconv.Itoa(profile_id) + ")\n"))
				return profile_id
			}
			return new_id
		}
	}
	return profile_id
}

func extract_args(og_args []string) []string {
	var args	[]string
	var val		string
	var pos		int
	var skip	bool

	skip = false
	for pos, val = range og_args {
		if val == "-p" {
			pos++
			skip = true
		} else if skip {
			skip = false
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
	str = string(str + bold + "kill" + reset + "	{TARGET}	Kill the indicated process\n")
	str = string(str + bold + "erase" + reset + "	{TARGET}	Erase the indicated process\n")
	str = string(str + bold + "russian" + reset + "			Russian roulet for your process\n")
	str = string(str + bold + "ps" + reset + "			List all the profiles and show their task id, within a profile id\n")
	str = string(str + bold + "ls" + reset + "	(ID)		List all the process within a profile id\n")
	str = string(str + bold + "ch" + reset + "	{ID}		Modify the working profile id\n")
	str = string(str + bold + "wichid" + reset + "			Print the current profile id in console\n")
	str = string(str + bold + "help" + reset + "			Show this info\n")
	str = string(str + bold + "-p" + reset + "			Set the profile id to the indicated number\n")
	str = string(str + "\nInformation in {} is a MUST and can't be skipped\n")
	str = string(str + "Information in () is a OPTIONAL and the id is the current id\n")
	str = string(str + "Current id is the the one modify by loading a configuration, using -p, or ch\n")
	return (str)
}
