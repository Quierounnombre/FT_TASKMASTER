package main

import "taskmaster/logger"

func main() {
	log, _ := logger.New("app.log")
	defer log.Close()
	
	log.Info("Started")
	log.Error("Test error")
}
