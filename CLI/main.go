package main

import (
	"net"
	"encoding/json"
)

const socket_path = "/run/taskmaster.sock"
const config_filename = "taskmaster.yaml"

type Cmd struct {
	Cmd			string		`json:"cmd"`
	Flags		[]string	`json:"flags"`
	Profile_id	int			`json:"profile_id"`
}

func main() {
	var sk		net.Conn
	var encoder	*json.Encoder
	
	sk = open_socket(socket_path)
	defer sk.Close()
	encoder = json.NewEncoder(sk)
	console_start(sk, encoder)
}
