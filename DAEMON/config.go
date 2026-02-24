package main

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

/*
Store the data from the yaml to launch a process
*/
type Process struct {
	Name            string            `yaml:"name"`            //THE PROGRAM
	Cmd             string            `yaml:"cmd"`             //cmd to launch the program
	Restart         string            `yaml:"restart"`         //always / never / on_error
	Stop_signal     string            `yaml:"stop_signal"`     //Signal used to gracefully stop ?minor bonus
	Work_dir        string            `yaml:"work_dir"`        //Working directory for the program
	Stdout          string            `yaml:"stdout"`          //Change the stdout to selected
	Stderr          string            `yaml:"stderr"`          //Change the stderr to selected
	Env             map[string]string `yaml:"env"`             //Enviroment variables for the program
	Restart_atempts int               `yaml:"restart_atemps"`  //Number of time would try to restart the program
	Expected_exit   []int             `yaml:"expected_exit"`   //Expected exit code
	Launch_wait     time.Duration     `yaml:"launch_wait"`     //Time until program launch is consider succesfull
	Kill_wait       time.Duration     `yaml:"kill_wait"`       //Waittime for killing the program
	Start_at_launch bool              `yaml:"start_at_launch"` //Start this program at launch or not
	Umask           *int              `yaml:"umask"`           //Umask to restrict permissions
	Num_procs       int               `yaml:"num_procs"`       //Number of this process to launch
}

/*
Store all the data from the yaml
*/
type File_Config struct {
	Process []Process `yaml:"process"`
	Path    string
}

//------------------------------------------------------------------------------------------------------------------------------------------

func get_file_content(name string) []byte {
	var content []byte
	var err error

	content, err = os.ReadFile(name)
	if err != nil {
		fmt.Println(err)
	}
	return content
}

func extract_file_content(raw_yaml []byte) *File_Config {
	var config File_Config
	var err error

	err = yaml.Unmarshal(raw_yaml, &config)
	if err != nil {
		fmt.Println(err)
	}
	return (&config)
}

// Change empty values for defaults
func set_config_defaults(config *File_Config) {
	for index := range config.Process {
		p := &config.Process[index]
		if p.Kill_wait.Seconds() < float64(time.Second) {
			p.Kill_wait = time.Second
		}
		if p.Launch_wait.Seconds() < float64(time.Second) {
			p.Launch_wait = time.Second
		}
		if len(p.Expected_exit) == 0 {
			p.Expected_exit = append(p.Expected_exit, 0)
		}
		if p.Restart == "" {
			p.Restart = "never"
		}
		if p.Umask == nil {
			defaultUmask := 22
			p.Umask = &defaultUmask // Default umask: owner rwx, group/others r-x
		}
		if p.Num_procs <= 0 {
			p.Num_procs = 1
		}
	}
}

func check_file_existance(path string) bool {
	var err error

	_, err = os.Stat(path)
	if err != nil {
		fmt.Println("Check for file existance in:", err)
		return (false)
	}
	return (true)
}

/*
Extract and process the file with that name, and returns a pointer s_config
with the data, only accepts .yaml In case of error, it exits
*/
func get_config_from_file_name(name string) *File_Config {
	var raw_yaml []byte
	var config *File_Config

	if check_file_existance(name) {
		raw_yaml = get_file_content(name)
		config = extract_file_content(raw_yaml)
		set_config_defaults(config)
		config.Path = name
	}
	return (config)
}

// JUST PRINTS
func PrintProcesStruct(p Process) {
	fmt.Println("=== Process Information ===")
	fmt.Printf("Name:             %s\n", p.Name)
	fmt.Printf("Command:          %s\n", p.Cmd)
	fmt.Printf("Restart Policy:   %s\n", p.Restart)
	fmt.Printf("Stop Signal:      %s\n", p.Stop_signal)
	fmt.Printf("Working Directory:%s\n", p.Work_dir)
	fmt.Printf("Stdout Path:      %s\n", p.Stdout)
	fmt.Printf("Stderr Path:      %s\n", p.Stderr)
	fmt.Printf("Environment Vars:\n")
	for key, value := range p.Env {
		fmt.Printf("  %s = %s\n", key, value)
	}
	fmt.Printf("Restart Attempts: %d\n", p.Restart_atempts)
	fmt.Printf("Expected Exits:   %v\n", p.Expected_exit)
	fmt.Printf("Launch Wait:      %s\n", p.Launch_wait)
	fmt.Printf("Kill Wait:        %s\n", p.Kill_wait)
	fmt.Printf("Start at Launch:  %t\n", p.Start_at_launch)
	fmt.Printf("Umask:             %d\n", p.Umask)
	fmt.Printf("Num_process:      %d\n", p.Num_procs)
	fmt.Println("============================")
}

// JUST PRINTS
func PrintFile_ConfigStruct(c File_Config) {
	fmt.Println("===     CONFIGURATION     ===")
	for _, element := range c.Process {
		PrintProcesStruct(element)
	}
	fmt.Println("=============================")
}
