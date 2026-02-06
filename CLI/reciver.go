package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

func get_id(json *map[string]interface{}) int {
	var id int
	var ok bool
	var idFloat float64

	idFloat, ok = (*json)["id"].(float64)
	if ok {
		id = int(idFloat)
	} else {
		id = -1
	}
	return id
}

func recive_load(json *map[string]interface{}, rl *readline.Instance, profile_id *int) {
	var flag string
	var id int
	var ok bool

	flag, ok = (*json)["flags"].(string)
	if !ok {
		flag = "ERROR MISSING CONTENT"
	}
	id = get_id(json)
	if id != -1 {
		*profile_id = id
		rl.Write([]byte("Loaded " + flag + "with id:" + strconv.Itoa(id) + "\n"))
	} else {
		rl.Write([]byte("Profile couldn't be loaded\n"))
	}
}

func recive_reload(json *map[string]interface{}, rl *readline.Instance) {
	var ok bool
	var flag string
	var id int

	flag, ok = (*json)["flags"].(string)
	if !ok {
		flag = "ERROR MISSING CONTENT"
	}
	id = get_id(json)
	if id != -1 {
		rl.Write([]byte("Stopped" + flag + "with id:" + strconv.Itoa(id) + "\n"))
		rl.Write([]byte("Loaded" + flag + "with id:" + strconv.Itoa(id) + "\n"))
	} else {
		rl.Write([]byte("Profile couldn't be reloaded\n"))
	}
}

func recive_stop(json *map[string]interface{}, rl *readline.Instance) {
	var id int

	id = get_id(json)
	if id != -1 {
		rl.Write([]byte("Stopped process with id:" + strconv.Itoa(id) + "\n"))
	} else {
		rl.Write([]byte("Process dosen't exist\n"))
	}
}

func recive_start(json *map[string]interface{}, rl *readline.Instance) {
	var id int

	id = get_id(json)
	if id != -1 {
		rl.Write([]byte("Started process with process id: " + strconv.Itoa(id) + "\n"))
	}
	if id != -1 {
		rl.Write([]byte("Started process with process id: " + strconv.Itoa(id) + "\n"))
	} else {
		rl.Write([]byte("Process dosen't exist\n"))
	}
}

func recive_restart(json *map[string]interface{}, rl *readline.Instance) {
	var id int

	id = get_id(json)
	if id != -1 {
		rl.Write([]byte("Stopped process with process id: " + strconv.Itoa(id) + "\n"))
		rl.Write([]byte("Started process with process id: " + strconv.Itoa(id) + "\n"))
	} else {
		rl.Write([]byte("Process dosen't exist\n"))
	}
}

func recive_describe(json *map[string]interface{}, rl *readline.Instance) {
	var fields []string
	var add func(string)
	var b strings.Builder
	var key string
	var value interface{}
	var env []interface{}

	add = func(key string) {
		b.WriteString(fmt.Sprintf("%-17s: %v\n", key, (*json)[key]))
	}
	fields = []string{
		"ID",
		"Name",
		"Status",
		"ExitCode",
		"RestartCount",
		"MaxRestarts",
		"StartTime",
		"WorkingDir",
		"Env",
		"ExpectedExitCodes",
		"Umask",
		"restartPolicy",
		"launchWait",
		"Cmd",
		"StdoutWriter",
		"StderrWriter",
	}
	for _, key = range fields {
		if key == "Env" {
			b.WriteString(fmt.Sprintf("%-17s:\n", key))
			env, _ = (*json)[key].([]interface{})
			for _, value = range env {
				b.WriteString(fmt.Sprintf("  - %v\n", value))
			}
		} else {
			add(key)
		}
	}
	rl.Write([]byte(b.String()))
}

func recive_error(json *map[string]interface{}, rl *readline.Instance) {
	var ok bool
	var flag string

	flag, ok = (*json)["flags"].(string)
	if !ok {
		flag = "ERROR MISSING CONTENT"
	}
	rl.Write([]byte("Error: " + flag + "\n"))
}

func recive_ps(json *map[string]interface{}, rl *readline.Instance) {
	var ok bool
	var path string
	var key string
	var keyFloat float64
	var prcs map[string]interface{}
	var proc_lst []interface{}
	var obj interface{}

	rl.Write([]byte("+------+------------------------------+\n"))
	rl.Write([]byte("| ID   | path                         |\n"))
	rl.Write([]byte("+------+------------------------------+\n"))
	proc_lst, ok = (*json)["profiles"].([]interface{})
	if !ok {
		proc_lst = nil
	}
	for _, obj = range proc_lst {
		prcs, ok = obj.(map[string]interface{})
		if !ok {
			prcs = nil
		}
		keyFloat, ok = prcs["profileID"].(float64)
		if ok {
			key = strconv.Itoa(int(keyFloat))
		} else {
			key = "-1"
		}
		path, ok = prcs["filePath"].(string)
		if !ok {
			path = "ERROR NO PATH"
		}
		rl.Write([]byte(
			fmt.Sprintf("| %-4s | %-28s |\n", key, path),
		))
	}
	rl.Write([]byte("+------+------------------------------+\n"))
}

func recive_ls(json *map[string]interface{}, rl *readline.Instance) {
	var ok				bool
	var proc_lst		[]interface{}
	var proc			map[string]interface{}
	var id				int
	var id_float		float64
	var name			string
	var status			string
	var ts				string
	var obj				interface{}

	rl.Write([]byte("+------+------------------+----------+-----------------------|\n"))
	rl.Write([]byte("|  ID  | Name             |  Status  | Timestamp             |\n"))
	rl.Write([]byte("+------+------------------+----------+-----------------------|\n"))
	proc_lst, ok = (*json)["procs"].([]interface{})
	if !ok {
		proc_lst = nil
	}
	for _, obj = range proc_lst {
		proc, ok = obj.(map[string]interface{})
		if !ok {
			proc = nil
		}
		id_float, ok = proc["taskID"].(float64)
		id = int(id_float)
		if !ok {
			id = -1
		}
		name, ok = proc["name"].(string)
		if !ok {
			name = "Null"
		}
		status, ok = proc["status"].(string)
		if !ok {
			status = "Null"
		}
		ts, ok = proc["timestamp"].(string)
		if !ok {
			ts = "Null"
		}
		rl.Write([]byte(fmt.Sprintf("| %-4s | %-16s | %-8s | %-21s |\n", strconv.Itoa(id), name, status, ts)))
	}
	rl.Write([]byte("+------+------------------+----------+-----------------------|\n"))
}

func recive_kill(json *map[string]interface{}, rl *readline.Instance) {
	var id int

	id = get_id(json)
	if id != -1 {
		rl.Write([]byte("Killed process with id: " + strconv.Itoa(id) + "\n"))
	} else {
		rl.Write([]byte("Process dosen't exists\n"))
	}
}

func recive_ch(json *map[string]interface{}, rl *readline.Instance, profile_id *int) {
	var id int

	id = get_id(json)
	if id != -1 {
		*profile_id = id
	}
	if id != -1 {
		rl.Write([]byte("Switched to id: " + strconv.Itoa(id) + "\n"))
	} else {
		rl.Write([]byte("Profile dosen't exists\n"))
	}
}

func reciver(json *map[string]interface{}, rl *readline.Instance, profile_id *int) {
	var ok bool
	var cmd string

	cmd, ok = (*json)["cmd"].(string)
	fmt.Println("CMD: ", cmd)
	fmt.Println("JSON: ", json)
	if !ok {
		rl.Write([]byte("ERROR CMD NOT FOUND"))
		return
	}
	switch cmd {
	case "load":
		recive_load(json, rl, profile_id)
	case "reload":
		recive_reload(json, rl)
	case "stop":
		recive_stop(json, rl)
	case "start":
		recive_start(json, rl)
	case "restart":
		recive_restart(json, rl)
	case "describe":
		recive_describe(json, rl)
	case "error":
		recive_error(json, rl)
	case "ps":
		recive_ps(json, rl)
	case "ls":
		recive_ls(json, rl)
	case "kill":
		recive_kill(json, rl)
	case "ch":
		recive_ch(json, rl, profile_id)
	default:
		rl.Write([]byte("DEFAULT"))
	}
}
