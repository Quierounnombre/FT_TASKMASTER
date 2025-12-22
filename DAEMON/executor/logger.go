package executor

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func New(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

func (l *Logger) Info(msg string) {
	l.write("INFO", msg)
}

func (l *Logger) Error(msg string) {
	l.write("ERROR", msg)
}

func (l *Logger) write(level, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, msg)
	l.file.WriteString(line)
}

func (l *Logger) Close() error {
	return l.file.Close()
}
