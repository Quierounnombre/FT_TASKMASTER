package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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

func format_duration(s string) string {
	if s == "" {
		return "null"
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return s
	}
	totalMs := d.Milliseconds()
	if totalMs < 1000 {
		return fmt.Sprintf("%dms", totalMs)
	}
	h := int(totalMs / 3600000)
	m := int((totalMs % 3600000) / 60000)
	sec := int((totalMs % 60000) / 1000)
	ms := int(totalMs % 1000)

	var b strings.Builder
	if h > 0 {
		fmt.Fprintf(&b, "%dh %dm %ds", h, m, sec)
	} else if m > 0 {
		fmt.Fprintf(&b, "%dm %ds", m, sec)
	} else {
		fmt.Fprintf(&b, "%ds", sec)
	}
	if ms > 0 {
		fmt.Fprintf(&b, " %dms", ms)
	}
	return b.String()
}

func enforce_max_size(str string, size int) string {
	var chars []rune
	var lenght int

	chars = []rune(str)
	lenght = len(chars)
	if lenght > size {
		return string(chars[lenght-size:])
	}
	return string(chars)
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
	var task		map[string]interface{}
	var ok			bool
	var fields		[]string
	var add			func(string, string)
	var b			strings.Builder
	var key			string
	var label		string
	var i			int
	var tmp			float64
	var value		interface{}
	var env			[]interface{}

	// The daemon wraps task info under a "task" key
	task, ok = (*json)["task"].(map[string]interface{})
	if !ok {
		rl.Write([]byte("ERROR: describe response missing 'task' field\n"))
		return
	}

	add = func(label string, key string) {
		b.WriteString(fmt.Sprintf("%-17s: %v\n", label, task[key]))
	}

	fields = []string{
		"id", "name", "status", "exitCode", "restartCount",
		"maxRestarts", "startTime", "workingDir", "env",
		"expectedExitCodes", "umask", "restartPolicy", "cmd",
	}
	labels := map[string]string{
		"id":                "ID",
		"name":              "Name",
		"status":            "Status",
		"exitCode":          "ExitCode",
		"restartCount":      "RestartCount",
		"maxRestarts":       "MaxRestarts",
		"startTime":         "StartTime",
		"workingDir":        "WorkingDir",
		"env":               "Env",
		"expectedExitCodes": "ExpectedExitCodes",
		"umask":             "Umask",
		"restartPolicy":     "RestartPolicy",
		"cmd":               "Cmd",
	}

	for _, key = range fields {
		label = labels[key]
		if key == "env" {
			env, _ = task[key].([]interface{})
			for i, value = range env {
				if i == 0 {
					b.WriteString(fmt.Sprintf("%-17s: - %v\n", label, value))
				} else {
					b.WriteString(fmt.Sprintf("%-19s- %v\n", "", value))
				}
			}
			if len(env) == 0 {
				b.WriteString(fmt.Sprintf("%-17s:\n", label))
			}
		} else if key =="umask" {
			tmp = task[key].(float64)
			i = int(tmp)
			b.WriteString(fmt.Sprintf("%-17s: %04o\n", label, i))
		} else {
			add(label, key)
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
	if len(os.Args) > start_shell {
		rl.Close()
		os.Exit(1)
	}
}

func recive_ps(json *map[string]interface{}, rl *readline.Instance) {
	var ok bool
	var path string
	var key string
	var keyFloat float64
	var prcs map[string]interface{}
	var proc_lst []interface{}
	var obj interface{}

	rl.Write([]byte("+------+-----------------------------------------------------+\n"))
	rl.Write([]byte("| ID   | path                                                |\n"))
	rl.Write([]byte("+------+-----------------------------------------------------+\n"))
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
		path = enforce_max_size(path, 51)
		rl.Write([]byte(
			fmt.Sprintf("| %-4s | %-51s |\n", key, path),
		))
	}
	rl.Write([]byte("+------+-----------------------------------------------------+\n"))
}

func recive_ls(json *map[string]interface{}, rl *readline.Instance) {
	var ok bool
	var proc_lst []interface{}
	var proc map[string]interface{}
	var id int
	var id_float float64
	var name string
	var status string
	var ts string
	var obj interface{}

	rl.Write([]byte("+------+--------------------------------+--------------+-----------------------|\n"))
	rl.Write([]byte("|  ID  | Name                           |  Status      | Time running          |\n"))
	rl.Write([]byte("+------+--------------------------------+--------------+-----------------------|\n"))
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
		ts, ok = proc["timeRunning"].(string)
		if !ok {
			ts = "Null"
		} else {
			ts = format_duration(ts)
		}
		name = enforce_max_size(name, 30)
		status = enforce_max_size(status, 12)
		ts = enforce_max_size(ts, 21)
		rl.Write([]byte(fmt.Sprintf("| %-4s | %-30s | %-12s | %21s |\n", strconv.Itoa(id), name, status, ts)))
	}
	rl.Write([]byte("+------+--------------------------------+--------------+-----------------------|\n"))
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
		rl.Write([]byte("Switched to id: " + strconv.Itoa(id) + "\n"))
	} else {
		rl.Write([]byte("Profile dosen't exists\n"))
	}
}

func recive_unload(json *map[string]interface{}, rl *readline.Instance, profile_id *int) {
	var id int

	id = get_id(json)
	if id != -1 {
		*profile_id = 0
		rl.Write([]byte("Profile unloadd\n"))
	} else {
		rl.Write([]byte("Error erasing profile"))
	}
}

func recive_russian(json *map[string]interface{}, rl *readline.Instance) {
	var unlucky string
	var ok bool

	unlucky, ok = (*json)["unlucky"].(string)
	if !ok {
		rl.Write([]byte("Need more players"))
		return
	}
	rl.Write([]byte("Our winner winner chiken dinner is: " + unlucky + "\n"))
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
	case "unload":
		recive_unload(json, rl, profile_id)
	case "russian":
		recive_russian(json, rl)
	default:
		rl.Write([]byte("DEFAULT"))
	}
}
