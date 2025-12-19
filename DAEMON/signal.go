package main

import (
	"os"
	"os/signal"
	"syscall"
	"taskmaster-daemon/executor"
)

/*
Create a channel and sets the flags and signals to listen
*/
func set_channel_for_signals() chan os.Signal {
	var channel chan os.Signal

	channel = make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGHUP)
	//ADD more sigint
	return (channel)
}

//SHOULD WE DO A RELOAD OF EVERY PROFILE?
func handle_signals(sig os.Signal, config []File_Config, manager *executor.Manager) {
	switch sig {
	case syscall.SIGHUP:
		var cmd Cmd

		cmd.base = "reload"
		//cmd.Execute(config, manager)
	default:
		os.Exit(1)
	}
}
