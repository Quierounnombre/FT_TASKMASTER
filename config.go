package main

import (
	"github.com/goccy/go-yaml"
	"fmt"
	"os"
)

/*
Store the data from the yaml to launch a process
*/
type Process struct {
	Name	string  `yaml:"name"`
	Args	string	`yaml:"args"`
}

/*
Store all the data from the yaml
*/
type Config struct {
	Process	[]Process `yaml:"process"`
	channel chan os.Signal
}

func	get_file_content(name string) []byte {
	var content []byte
	var err error
	
	content, err = os.ReadFile(name)
	if (err != nil) {
		fmt.Println(err)
		os.Exit(1)
	}
	return content
}

func extract_file_content(raw_yaml []byte) *Config {
	var config	Config
	var err		error

	err = yaml.Unmarshal(raw_yaml, &config)
	if (err != nil) {
		fmt.Println(err)
		os.Exit(1)
	}
	return (&config)
}

/*
Extract and process the file with that name, and returns a pointer s_config with the data, only accepts .yaml
In case of error, it exits
*/
func get_config_from_file_name(name string) *Config {
	var	raw_yaml	[]byte
	var config		*Config

	raw_yaml = get_file_content(name)
	config = extract_file_content(raw_yaml)
	return (config)
}