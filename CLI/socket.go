package main

import (
	"net"
	"os"
	"fmt"
	"bufio"
	"github.com/chzyer/readline"
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

func send_data(sk net.Conn, str string) {
	var err		error
	var bytes	[]byte	

	bytes = append([]byte(str), '\n')
	_, err = sk.Write(bytes)
	if (err != nil) {
		fmt.Println("Error socket not working")
		os.Exit(1)
	}
}

func recive_data(sk net.Conn, rl *readline.Instance) {
	var reader		*bufio.Reader
	var msg			string
	var err			error

	reader = bufio.NewReader(sk)
	for (true) {
		msg, err = reader.ReadString('\n')
		if (err != nil) {
			fmt.Println("ERROR_RECIVING_DATA")
			break
		}
		rl.Write([]byte(msg))
		rl.Refresh()
	}
}