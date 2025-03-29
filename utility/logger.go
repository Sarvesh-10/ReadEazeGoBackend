package utility

import "log"

type Logger struct {}


func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Info(message string,args ...interface{}) {
	log.Printf("\n[INFO] " + message, args...)
}

func (l *Logger) Error(message string,args ...interface{}) {
	log.Printf("\n[ERROR] " + message, args...)
}



