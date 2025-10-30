package main

import (
	"fmt"
	"os"
)

const socket_path = "/run/taskmaster.sock"

func main() {
	var sock_config	Sock_Config

	sock_config.sig_ch = set_channel_for_signals()
	sock_config.cli_ch = socket_wrapper(socket_path, &sock_config)

	//PrintConfigStruct(*config)
	
	loop(&sock_config)
}

/*
Main loop with signal support
*/
func loop(sock_config *Sock_Config) {
	var file_config	*File_Config
	var signal 		os.Signal
	var msg			string
	
	for (true) {
		select {
		case signal = <- sock_config.sig_ch:
			fmt.Println("SIGNAL", "=", signal)
			os.Exit(1)
		case msg = <- sock_config.cli_ch:
			fmt.Println("MSG", " = ", msg)
			file_config = get_config_from_file_name(msg, file_config)
			broadcast_data(sock_config.cons, msg)
		default:
			//
		}
	}
}