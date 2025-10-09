package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	var config *Config

	config = get_config_from_file_name("example.yaml")
	config.channel = set_channel_for_signals()

	fmt.Println(config.Process[0].Name)
	fmt.Println(config.Process[1].Args)
	fmt.Println(config.Process[1].Name)
	fmt.Println(config.Process[0].Args)
	
	loop(config)
}

/*
Main loop with signal support
*/
func loop(config *Config) {
	var signal os.Signal

	for (true) {
		select {
		case signal = <- config.channel:
			fmt.Println("SIGNAL", "=", signal)
			os.Exit(1)
		default:
			fmt.Println("WAITING")
			time.Sleep(time.Second)
		}
	}
}