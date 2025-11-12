package main

import (
	"fmt"
	"net"
	"os"
)

const socket_path = "/run/taskmaster.sock"
const config_filename = "taskmaster.yaml"

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
	var	sk		net.Conn
	var path	string
	
	path  = get_yaml_path()
	check_file_existance(path)
	sk = open_socket(socket_path)
	defer sk.Close()
	send_data(sk, path)
	console_start(sk)
}
