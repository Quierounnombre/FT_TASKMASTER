package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

//Small funtion for creating the socket
func create_socket(socket_path string) net.Listener {
	var sk	net.Listener
	var err	error

	os.Remove(socket_path)
	sk, err = net.Listen("unix", socket_path)
	if (err != nil) {
		fmt.Println("ERROR_CREATING_SCOKET")
		os.Exit(1)
	}
	return (sk)
}

func handle_connection(sk net.Listener, ch chan string) {
	var con		net.Conn
	var err		error
	
	defer sk.Close()
	for (true) {
		con, err = sk.Accept()
		if (err != nil) {
			fmt.Println(err)
			os.Exit(1)
		}
		go handle_client(con, ch)
	}
}

func handle_client(conn net.Conn, ch chan string) {
	var reader		*bufio.Reader
	var msg			string
	var err			error

	reader = bufio.NewReader(conn)
	for (true) {
		msg, err = reader.ReadString('\n')
		if (err != nil) {
			conn.Close()
			if (err == io.EOF) {
				fmt.Println("DISCONNECTION")
				break
			}
			fmt.Println("ERROR_RECIVING_DATA")
			break
		}
		ch <- msg
	}
}

//Generates a go channel and starts the subrutine for sending data through the channel
func socket_wrapper(socket_path string) chan string {
	var ch chan string
	var sk net.Listener

	sk = create_socket(socket_path)
	ch = make(chan string)
	go handle_connection(sk, ch)
	return  (ch)
}