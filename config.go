package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

/*
Expand this struct as needed
*/
type Config struct {
	To		string	`xml:"to"`
	From	string	`xml:"from"`
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

func extract_file_content(raw_xml []byte) *Config {
	var config	Config
	var err		error

	err = xml.Unmarshal(raw_xml, &config)
	if (err != nil) {
		fmt.Println(err)
		os.Exit(1)
	}
	return (&config)
}

/*
Extract and process the file with that name, and returns a pointer s_config with the data, only accepts .xml
In case of error, it exits
*/
func get_config_from_file_name(name string) *Config {
	var	raw_xml	[]byte
	var config	*Config

	raw_xml = get_file_content(name)
	config = extract_file_content(raw_xml)
	return (config)
}