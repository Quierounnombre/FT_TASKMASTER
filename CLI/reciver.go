import {
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

//NEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEED CONTENT
func recive_describe(json *map[string]interface{}) {
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

//NEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEED CONTENT
func recive_ps(json *map[string]interface{}) {
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

//NEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEED CONTENT
func recive_ls(json *map[string]interface{}) {
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

//NEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEED CONTENT
func recive_help(json *map[string]interface{}) {
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
	case "help":
		recive_help(json)
	}
}
