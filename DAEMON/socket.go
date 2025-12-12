package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"encoding/json"
)

type Sock_Config struct {
	sig_ch chan os.Signal
	cli_ch chan Msg
	cons   []net.Conn
}

// Small funtion for creating the socket
func create_socket(socket_path string) net.Listener {
	var sk net.Listener
	var err error

	os.Remove(socket_path)
	sk, err = net.Listen("unix", socket_path)
	if err != nil {
		fmt.Println(err)
		fmt.Println("RUN WITH SUDO")
		os.Exit(1)
	}
	return (sk)
}

// Manage the creation of the connection with client
func handle_connection(sk net.Listener, ch chan Msg, config *Sock_Config) {
	var con net.Conn
	var err error

	defer sk.Close()
	for true {
		con, err = sk.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		config.cons = append(config.cons, con)
		go handle_client(con, ch, config)
	}
}

// Clean the list of clients
func remove_closed_client(target_conn net.Conn, config *Sock_Config) {
	var cons []net.Conn

	for _, element := range config.cons {
		if element != target_conn {
			cons = append(cons, element)
		}
	}
	config.cons = cons
}

// HANDLE THE CLIENT RECIVING DATA
func handle_client(conn net.Conn, ch chan Msg, config *Sock_Config) {
	var decoder *json.Decoder
	var msg Msg
	var err error

	msg.author = conn
	msg.encoder = json.NewEncoder(conn)
	decoder = json.NewDecoder(conn)
	for true {
		err = decoder.Decode(&msg.content)
		if err != nil {
			remove_closed_client(conn, config)
			conn.Close()
			if err == io.EOF {
				fmt.Println("DISCONNECTION")
				break
			}
			fmt.Println("ERROR_RECIVING_DATA")
			break
		}
		ch <- msg
	}
}

// SEND DATA TO ALL THE CLIENTS
func broadcast_data(connections []net.Conn, str string) {
	var err error
	var bytes []byte

	bytes = append([]byte(str), '\n')
	for _, conn := range connections {
		_, err = conn.Write(bytes)
		if err != nil {
			fmt.Println("Error socket not working")
			fmt.Println(err)
			fmt.Println("Target conn -> ", conn)
		}
	}
}

// Generates a go channel and starts the subrutine for sending data through the channel
func socket_wrapper(socket_path string, config *Sock_Config) chan Msg {
	var ch chan Msg
	var sk net.Listener

	sk = create_socket(socket_path)
	ch = make(chan Msg)
	go handle_connection(sk, ch, config)
	return (ch)
}
