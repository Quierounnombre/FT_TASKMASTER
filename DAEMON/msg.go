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
	var ok		bool
	var value	string

	value, ok = m.content["cmd"].(string)
	if (ok) {
		return (value)
	}
	return ("")
}

func (m *Msg) get_flags() []string {
	var ok		bool
	var raw		[]interface{}
	var flags	[]string

	raw, ok = m.content["flags"].([]interface{})
	if (ok) {
		for _, value := range raw {
			flags = append(flags, value.(string))
		}
		return (flags)
	}
	return (nil);
}

func (m *Msg) get_profile_id() int {
	var ok			bool
	var profile_id	int

	profile_id, ok = m.content["profile_id"].(int)
	if (ok) {
		return (value)
	}
	return (-1)
}

func (m *Msg) add_payload(key string, value interface{}) {
	m.content[key] = value
}
