package main

import (
	"net"
	"fmt"
	"strings"
)

type	Msg struct {
	author		net.Conn
	content		string
}

func (m *Msg)clean_content() {
	m.content = strings.Trim(m.content, "\n") // DATA CLEAN UP
}

func (m *Msg)reply(reply string) {
	var bytes	[]byte
	var	err		error

	bytes = append([]byte(reply), '\n')
	_, err = m.author.Write(bytes)
	if (err != nil) {
		fmt.Println("Error socket not working")
		fmt.Println(err)
		fmt.Println("Target conn -> ", m.author)
	}
}

func (m *Msg)print_msg() {
	fmt.Printf("{\n auth = %v\n content = %s\n}\n", m.author, m.content)
}