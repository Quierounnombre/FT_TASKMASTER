package main

import (
	"fmt"
)

func main() {
	var config *Config

	config = get_config_from_file_name("example.yaml")

	fmt.Println(config.Process[0].Name)
	fmt.Println(config.Process[1].Args)
	fmt.Println(config.Process[1].Name)
	fmt.Println(config.Process[0].Args)
}
