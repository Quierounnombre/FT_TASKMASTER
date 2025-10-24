package main

import (
	"fmt"
	"os"
	"net"
	"github.com/chzyer/readline"
)

//Sets the configuration for the console
func	set_config() *readline.Config {
	var config_rl readline.Config

	config_rl.HistoryFile = ".history"
	config_rl.Prompt = "<>< "
	config_rl.HistoryLimit = 100
	return (&config_rl)
}

//Sets up the channel and starts the go subrutine
func	console_start(sk net.Conn) {
	var rl			*readline.Instance
	var config_rl	*readline.Config
	var err			error

	config_rl = set_config()
	rl, err = readline.NewEx(config_rl)
	if (err != nil) {
		os.Exit(1)
	}
	console(rl, sk)
}

//ASYNCRONOUS FUNC call with go
func	console(rl *readline.Instance, sk net.Conn) {
	var err		error
	var line	string
	
	for (true) {
		line, err = rl.Readline()
		if (err != nil) {
			os.Exit(1)
		}
		fmt.Println("Has escrito -> ", line)
		send_data(sk, line)
	}
}