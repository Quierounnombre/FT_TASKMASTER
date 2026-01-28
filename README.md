## How to develop

Start up the environment:
```
docker compose up -d
```

Then access the container:
```
docker compose exec taskmaster-dev bash
```

## Structure of the project

Taskmaster is going to work with a client & server way. There are two directories corresponding to the CLI (client) and DAEMON (server).

### CLI structure

IMPORTANT: DAEMON MUST BE RUNNING

How to use:

Open terminal:
Execute the cmd(./CLI) without any arguments this will open a terminal for you to work
with our cmds.

CLI:
Execute the cmd(./CLI) with arguments, this will send the cmd, and wait for the response, then return control
to the user.

Internals:
- `socket.go` General networking.
- `reciver.go` Boilerplate for output from daemon.
- `console.go` Actual logic for the CLI, either send a single cmd or remains open with rl.

### Daemon structure

The main configuration of a server is done at the root of the directory:
- `cmd.go` has the inputs to the server. The function execute connects with the `manager.go`
- `config.go`
- `logger.go`
- `msg.go`
- `signal.go`
- `socket.go`

There is `executor` directory with the manager and executor:
- `manager.go` will manage profiles and the executor. It stores each file configuration as a profile. Each profile has a pointer to an Executor struct.
- `executor.go` has a struct Executor with a map of the struct Task, that stores each process configuration.
- TODO: a "vigilante" that will start, restart, kill,... if the configuration requiere.

> Manager is a midle point between Executor/Vigilante and the CLI request + upper daemon structure.
