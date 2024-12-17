package logger

import (
	"fmt"
	"os"
)

type fileLogger struct {
	file          *os.File
	isLogRequired bool
}

func NewFileLogger(filename string, isLogRequired bool) (Logger, func() error, error) {

	if !isLogRequired {
		return &fileLogger{file: nil, isLogRequired: false}, func() error { return nil }, nil
	}

	// Create or open the log file
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return &fileLogger{file: file, isLogRequired: true}, file.Close, nil
}

// Log an info message
func (f *fileLogger) Log(message string) {
	if !f.isLogRequired {
		return
	}
	f.file.Write([]byte(fmt.Sprintf("[MSG] %s \n", message)))
}
