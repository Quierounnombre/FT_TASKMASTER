package main

import (
	"os"
	"net"
	"github.com/chzyer/readline"
	"encoding/json"
)

const history_path = "/run/.history"

//Sets the configuration for the console
func	set_config() *readline.Config {
	var config_rl readline.Config

	config_rl.HistoryFile = history_path
	config_rl.Prompt = "<>< "
	config_rl.HistoryLimit = 100
	return (&config_rl)
}

//Sets up rl library and starts the console
func	console_start(sk net.Conn, encoder *json.Encoder) {
	var rl			*readline.Instance
	var config_rl	*readline.Config
	var err			error

	config_rl = set_config()
	rl, err = readline.NewEx(config_rl)
	if (err != nil) {
		os.Exit(1)
	}
	go recive_data(sk, rl)
	console(rl, encoder)
}

//Starts the console
func	console(rl *readline.Instance, encoder *json.Encoder) {
	var err		error
	var line	string
	var cmd		Cmd
	
	for (true) {
		line, err = rl.Readline()
		if (err != nil) {
			os.Exit(1)
		}
		cmd.cmd = line
		send_data(encoder, &cmd)
	}
}
