package main

import (
	"fmt"
	"os"
)

const socket_path = "/run/taskmaster.sock"

func main() {
	var config *Config

	config = get_config_from_file_name("example.yaml")
	config.sig_ch = set_channel_for_signals()
	config.cli_ch = socket_wrapper(socket_path)

	//PrintConfigStruct(*config)

	loop(config)
}

/*
Main loop with signal support
*/
func loop(config *Config) {
	var signal 	os.Signal
	var msg		string

	for (true) {
		select {
		case signal = <- config.sig_ch:
			fmt.Println("SIGNAL", "=", signal)
			os.Exit(1)
		case msg = <- config.cli_ch:
			fmt.Print("MSG", "=", msg)
		default:
			//
		}
	}
}