package main

import (
	"os"
	"fmt"
	"taskmaster-daemon/executor"
)

const socket_path = "/run/taskmaster.sock"

func main() {
	var sock_config Sock_Config

	sock_config.sig_ch = set_channel_for_signals()
	sock_config.cli_ch = socket_wrapper(socket_path, &sock_config)

	//PrintConfigStruct(*config)

	loop(&sock_config)
}

func PrintMap(m map[string]interface{}) {
	for k, v := range m {
		fmt.Printf("%s: %v\n", k, v)
	}
}

/*
Main loop with signal support
*/
func loop(sock_config *Sock_Config) {
	var file_config []File_Config
	var signal os.Signal
	var manager *executor.Manager
	var msg Msg
	var cmd Cmd

	manager = executor.NewManager()
	for true {
		select {
		case signal = <-sock_config.sig_ch:
			handle_signals(signal, file_config, manager)
		case msg = <-sock_config.cli_ch:
			cmd.empty_cmd()
			PrintMap(msg.content)
			cmd.Parse_cmd(&msg)
			cmd.Execute(file_config, manager, msg)
			msg.add_payload("response");
			msg.reply()
			//broadcast_data(sock_config.cons, msg)
		default:
			//
		}
	}
}
