package main

import (
	"fmt"
	"time"

	"taskmaster/executor"
	"taskmaster/logger"
)

func main() {
	log, _ := logger.New("app.log")
	defer log.Close() //This will ensure the log file is closed when main exits

	// Example log entries
	log.Info("Started")
	log.Error("Test error")

	// Executor example
	exec := executor.New()

	// Start a long-running task
	exec.Execute("task1", "/tmp/task1.log", "sleep", "5")

	// Check status
	time.Sleep(100 * time.Millisecond)
	status, _ := exec.GetStatus("task1")
	fmt.Println("Status:", status)

	// Stop the task
	exec.Stop("task1")

	time.Sleep(100 * time.Millisecond)
	status, _ = exec.GetStatus("task1")
	fmt.Println("Final status:", status)

	// Output:
	// Status: running
	// Final status: stopped
}
