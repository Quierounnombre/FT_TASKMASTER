package main

import (
	"net"
	"os"
	"fmt"
	"github.com/chzyer/readline"
	"encoding/json"
)

func open_socket(socket_path string) net.Conn {
	var sk net.Conn
	var err error

	sk, err = net.Dial("unix", socket_path)
	if (err != nil) {
		fmt.Println("Can't connect to socket, make sure the daemon is up and running")
		os.Exit(1)
	}
	return (sk)
}

func send_data(encoder *json.Encoder, cmd *Cmd) {
	var err		error

	err = encoder.Encode(cmd)
	if (err != nil) {
		fmt.Println("Error socket not working")
		os.Exit(1)
	}
}

func recive_data(sk net.Conn, rl *readline.Instance, profile_id *int, recived chan struct{}) {
	var decoder		*json.Decoder
	var msg			map[string]interface{}
	var err			error

	decoder = json.NewDecoder(sk)
	for (true) {
		err = decoder.Decode(&msg)
		if (err != nil) {
			fmt.Println("ERROR_RECIVING_DATA")
			break
		}
		reciver(&msg, rl, profile_id)
		rl.Refresh()
		if (len(os.Args) > start_shell) {
			os.Exit(0)
		}
		recived <- struct{}{}
	}
}

func PrintMapRL(m map[string]interface{}, rl *readline.Instance) {
    b, err := json.MarshalIndent(m, "", "  ")
    if err != nil {
        rl.Write([]byte("Error marshaling map: " + err.Error() + "\n"))
        return
    }
    rl.Write(append(b, '\n'))
}
