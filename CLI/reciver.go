import {
	"fmt"
	"github.com/chzyer/readline"
}

func recive_load(json *map[string]interface{}) {
	var flag	string
	var id		int
	var ok		bool

	flag, ok = json["flags"].(string)
	if (!ok) {
		flag = "ERROR MISSING CONTENT"
	}
	id, ok = json["id"].(int)
	if (!ok) {
		id = -1
	}
	rl.Write([]byte("Loaded " + flag + "with id:" + id + "\n"))
}

func recive_reload(json *map[string]interface{}) {
	var ok		bool
	var flag	string
	var id		int

	flag, ok = json["flags"].(string)
	if (!ok) {
		flag = "ERROR MISSING CONTENT"
	}
	id, ok = json["id"].(int)
	if (!ok) {
		id = -1
	}
	rl.Write([]byte("Stopped " + flag + "with id:" + id + "\n"))
	rl.Write([]byte("Loaded " + flag + "with id:" + id + "\n"))
}

func recive_stop(json *map[string]interface{}) {
	var ok		bool
	var id		int

	id, ok = json["id"].(int)
	if (!ok) {
		id = -1
	}
	rl.Write([]byte("Stopped " + flag + "with id:" + id + "\n"))
}

func recive_start(json *map[string]interface{}) {
	var ok		bool
	var flag	string
	var id		int

	flag, ok = json["flags"].(string)
	if (!ok) {
		flag = "ERROR MISSING CONTENT"
	}
	id, ok = json["id"].(int)
	if (!ok) {
		id = -1
	}
	rl.Write([]byte("Started process with process id: " + id + "\n"));
}

func recive_restart(json *map[string]interface{}) {
	var ok		bool
	var flag	string
	var id		int

	flag, ok = json["flags"].(string)
	if (!ok) {
		flag = "ERROR MISSING CONTENT"
	}
	id, ok = json["id"].(int)
	if (!ok) {
		id = -1
	}
	rl.Write([]byte("Stopped process with process id: " + id + "\n"));
	rl.Write([]byte("Started process with process id: " + id + "\n"));
}

func recive_describe(json *map[string]interface{}) {
	var fields	[]string
	var add		func(string)
	var b		strings.Builder
	var key		string
	var value	string
	var env		[]interface{}

	add = func(key string) {
		b.WriteString(fmt.Sprintf("%-17s: %v\n", key, m[key]))
	}
	fields = []string {
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
		if (key == "Env") {
			b.WriteString(fmt.Sprintf("%-17s:\n", key))
			env, _ = json[key].([]interface{})
			for _, value = range env {
				b.WriteString(fmt.Sprintf("  - %v\n", v)
			}
		} else {
			add(key)
		}
	}
	rl.Write([]byte(b.String()))
}

func recive_error(json *map[string]interface{}) {
	var ok		bool
	var flag	string
	var id		int

	flag, ok = json["flags"].(string)
	if (!ok) {
		flag = "ERROR MISSING CONTENT"
	}
	id, ok = json["id"].(int)
	if (!ok) {
		id = -1
	}
	rl.Write([]byte("Error: " + flag + "\n"));
}

func recive_ps(json *map[string]interface{}) {
	var ok		bool
	var path	string
	var key		string
	var prcs	map[string]interface{}
	var value	string

	prcs, ok = json["prcs"].(map[string]interface{})
	if (!ok) {
        prcs = nil
	}
	rl.Write([]byte("+------+------------------------------+\n"))
	rl.Write([]byte("| ID   | path                         |\n"))
	rl.Write([]byte("+------+------------------------------+\n"))
	for key, value = range prcs {
		path, ok = prcs["path"].(string)
		if (!ok) {
			value = "ERROR MISSING CONTENT"
		}
 		rl.Write([]byte(
			fmt.Sprintf("| %-4s | %-28s |\n", key, path),
		))
	}
	rl.Write([]byte("+------+------------------------------+\n"))
}

func recive_ls(json *map[string]interface{}) {
	var ok			bool
	var proc_lst	[]interface{}
	var proc		map[string]interface{}
	var id			int
	var name		string
	var status		string
	var ts			string

	rl.Write([]byte("+----+------------+--------+-----------------------|\n"))
	rl.Write([]byte("| ID | Name       | Status | Timestamp             |\n"))
	rl.Write([]byte("+----+------------+--------+-----------------------|\n"))
	proc_lst, ok = json["procs"].([]interface{})
	if (!ok) {
         proc_lst = nil
	}
	for _, obj = range proc_lst {
		proc, ok = obj.(map[string]interface{})
		if (!ok) {
			proc = nill
		}
		id, ok = proc["id"].(string)
		if !ok {
			id = "Null"
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
		rl.Write([]byte(fmt.Sprintf("| %-4s | %-12s | %-8s | %-24s |\n", id, name, status, ts)))
	}
	rl.Write([]byte("+----+------------+--------+-----------------------|\n"))
}

func recive_kill(json *map[string]interface{}) {
	var ok		bool
	var id		int

	id, ok = json["id"].(int)
	if (!ok) {
		id = -1
	}
	rl.Write([]byte("Killed process with id: " + id + "\n"));
}

//ADD KILL
func reciver(json *map[string]interface{}) {
	var ok		bool
	var cmd		string

	switch cmd {
	case "load":
		recive_load(json)
	case "reload":
		recive_reload(json)
	case "stop":
		recive_stop(json)
	case "start":
		recive_start(json)
	case "restart":
		recive_restart(json)
	case "describe":
		recive_descrive(json)
	case "error":
		recive_error(json)
	case "ps":
		recive_ps(json)
	case "ls":
		recive_ls(json)
	case "kill":
		recive_kill(json)
	}
}
