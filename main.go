package main

import (
	"fmt"
)

func main() {
	var config *Config

	config = get_config_from_file_name("example.xml")

	fmt.Println(config.To)
	fmt.Println(config.From)
}
