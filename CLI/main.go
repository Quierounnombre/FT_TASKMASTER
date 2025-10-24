package main

import (
	"net"
)

const socket_path = "/run/taskmaster.sock"

func main() {
	var	sk net.Conn

	sk = open_socket(socket_path)
	console_start(sk)
}
