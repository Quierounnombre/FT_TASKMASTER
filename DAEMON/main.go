package main

import (
	"os"
	"taskmaster-daemon/executor"
)

const socket_path = "/run/taskmaster.sock"

func main() {
	var sock_config Sock_Config
	manager := executor.NewManager()
	logger := manager.Logger()

	sock_config.sig_ch = set_channel_for_signals()
	sock_config.cli_ch = socket_wrapper(socket_path, &sock_config, logger)

	//PrintConfigStruct(*config)

	loop(&sock_config, manager)
}

/*
Main loop with signal support
*/
func loop(sock_config *Sock_Config, manager *executor.Manager) {
	var signal os.Signal
	var msg Msg
	var cmd Cmd

	for true {
		select {
		case signal = <-sock_config.sig_ch:
			handle_signals(signal)
		case msg = <-sock_config.cli_ch:
			cmd.empty_cmd()
			cmd.Parse_cmd(&msg)
			cmd.Execute(manager, &msg)
			msg.reply()
			//broadcast_data(sock_config.cons, msg)
		default:
			//
		}
	}
}
