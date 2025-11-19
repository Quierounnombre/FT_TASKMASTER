package main

import (
	"os"
)

const socket_path = "/run/taskmaster.sock"

func main() {
	var sock_config Sock_Config

	sock_config.sig_ch = set_channel_for_signals()
	sock_config.cli_ch = socket_wrapper(socket_path, &sock_config)

	//PrintConfigStruct(*config)

	loop(&sock_config)
}

/*
Main loop with signal support
*/
func loop(sock_config *Sock_Config) {
	var file_config []File_Config
	var signal os.Signal
	var msg Msg
	var cmd Cmd

	for true {
		select {
		case signal = <-sock_config.sig_ch:
			handle_signals(signal, file_config)
		case msg = <-sock_config.cli_ch:
			cmd.empty_cmd()
			msg.print_msg()
			cmd.Parse_cmd(msg.content)
			tmp := cmd.Execute(file_config)
			msg.reply(tmp)
			//broadcast_data(sock_config.cons, msg)
		default:
			//
		}
	}
}
