package main

import (
	"fmt"
	"os"
)

func main() {
	var config *Config

	config = get_config_from_file_name("example.yaml")
	config.sig_ch = set_channel_for_signals()
	config.input_ch = console_start()

	//PrintConfigStruct(*config)

	loop(config)
}

/*
Main loop with signal support
*/
func loop(config *Config) {
	var signal 	os.Signal
	var text	string

	for (true) {
		select {
		case signal = <- config.sig_ch:
			fmt.Println("SIGNAL", "=", signal)
			os.Exit(1)
		case text = <- config.input_ch:
			fmt.Println("ESCRIBISTE", text)
		default:
			//
		}
	}
}