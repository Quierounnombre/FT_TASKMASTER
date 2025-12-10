package main

import (
	"fmt"
	"net"
	"encoding/json"
)

type Msg struct {
	author  net.Conn
	encoder *json.Encoder
	content map[string]interface{}
}

func (m *Msg) clean_content() {
	m.content = make(map[string]interface{})
}

func (m *Msg) reply() {
	var err error

	err = m.encoder.Encode(m.content)
	if err != nil {
		fmt.Println("Error socket not working")
		fmt.Println(err)
		fmt.Println("Target conn -> ", m.author)
	}
}

func (m *Msg) get_cmd() string {
	return (m.content["cmd"].(string))
}

func (m *Msg) get_flags() []string {
	var raw		[]interface{}
	var flags	[]string

	raw = m.content["flags"].([]interface{})
	for _, value := range raw {
		flags = append(flags, value.(string))
	}
	return (flags)
}

func (m *Msg) add_payload(key string, value interface{}) {
	m.content[key] = value
}
