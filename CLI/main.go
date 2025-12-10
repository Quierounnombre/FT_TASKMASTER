package main

import (
	"fmt"
	"net"
	"os"
	"encoding/json"
)

const socket_path = "/run/taskmaster.sock"
const config_filename = "taskmaster.yaml"

type Cmd struct {
	Cmd		string		`json:"cmd"`
	Flags	[]string	`json:"flags"`
}

func get_yaml_path() string {
	var cwd			string
	var err			error

	cwd, err = os.Getwd()
	if (err != nil) {
		fmt.Println("Can't get CWD")
		os.Exit(1)
	}
	cwd = cwd + string(os.PathSeparator) + config_filename
	return (cwd)
}

func check_file_existance(path string) {
	var	err	error

	_, err = os.Stat(path)
	if (err != nil) {
		fmt.Println("Check for file existance in:", path)
		os.Exit(1)
	}
}

func main() {
	var sk		net.Conn
	var encoder	*json.Encoder
	var path	string
	var cmd		Cmd
	
	path  = get_yaml_path()
	check_file_existance(path)
	sk = open_socket(socket_path)
	defer sk.Close()
	encoder = json.NewEncoder(sk)
	cmd.Cmd = "load"
	cmd.Flags = append(cmd.Flags, path)
	send_data(encoder, &cmd)
	console_start(sk, encoder)
}
