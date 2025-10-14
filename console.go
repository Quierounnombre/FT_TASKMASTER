package main

import (
	"github.com/chzyer/readline"
	"os"
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
func	console_start() chan string {
	var rl			*readline.Instance
	var config_rl	*readline.Config
	var input_ch 	chan string
	var err			error

	config_rl = set_config()
	input_ch = make(chan string)
	rl, err = readline.NewEx(config_rl)
	if (err != nil) {
		os.Exit(1)
	}
	go console(input_ch, rl)
	return (input_ch)
}

//ASYNCRONOUS FUNC call with go
func	console(input_ch chan string, rl *readline.Instance) {
	var err		error
	var line	string
	
	for (true) {
		line, err = rl.Readline()
		if (err != nil) {
			os.Exit(1)
		}
		input_ch <- line
	}
}