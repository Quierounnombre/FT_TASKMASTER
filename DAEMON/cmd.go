package main

import (
	"fmt"
	"strings"
)

const bold = "\033[1m"
const reset = "\033[0m"

type Cmd struct {
	base  string
	flags []string
	err   bool
}

// Parses the cmd and do all the data cleanup
func (c *Cmd) Parse_cmd(content string) {
	var splited_content []string

	splited_content = strings.Split(content, " ")
	if len(splited_content) > 0 {
		c.base = splited_content[0]
		for _, element := range splited_content {
			if element != c.base {
				c.flags = append(c.flags, element)
			}
		}
		//NEED FOR FLAG_CHEKING IN THE FUTURE?
	} else {
		fmt.Println("EMPTY CONTENT")
	}
}

// EXECUTE THE CMD, IN THE FUTURE WILL MERGE WITH PABLO
func (c *Cmd) Execute(config []File_Config) string {
	switch c.base {
	case "load":
		tmp := get_config_from_file_name(c.flags[0])
		config = append(config, *tmp)
		PrintFile_ConfigStruct(config[0])
		return (string("Loaded " + c.flags[0]))
	case "reload":
		for index, element := range config {
			config[index] = *get_config_from_file_name(element.Path)
		}
		return ("Configurations reloaded")
	case "stop":
		//CHECK FOR AVIABLE PROCESS HERE?
		return (string("Stoped " + c.flags[0]))
	case "start":
		//CHECK FOR AVIABLE PROCESS HERE?
		return (string("Started " + c.flags[0]))
	case "restart":
		if len(c.flags) > 0 {
			return (string("Restarted " + c.flags[0]))
		}
		return (string("Restarting all programs"))
	case "ls":
		return ("HOLA")
	case "help":
		return (cmd_help())
	}
	return ("Wrong cmd")
}

func (c *Cmd) empty_cmd() {
	c.base = ""
	c.flags = nil
}

// Returns a list with all the current working cmds
func cmd_help() string {
	var str string

	str = string(str + bold + "load" + reset + "	{PATH} load the taskmaster.yaml file in the current dir\n")
	str = string(str + bold + "reload" + reset + "	reload all the configuration files executing in the taskmaster\n")
	str = string(str + bold + "stop" + reset + "	{TARGET} stop the target process\n")
	str = string(str + bold + "start" + reset + "	{TARGET} start the target process\n")
	str = string(str + bold + "restart" + reset + "	(TARGET) reset the target process, or all of them in case there are no target\n")
	str = string(str + bold + "help" + reset + "	Show this info")
	return (str)
}
