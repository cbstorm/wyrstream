package logger

import (
	"fmt"
	"log"
)

var UnexpectedErrLogger = NewLogger("UNEXPECTED")

type Logger struct {
	Namespace string
}

func NewLogger(ns string) *Logger {
	return &Logger{
		Namespace: ns,
	}
}

func (l *Logger) Child(ns string) *Logger {
	return &Logger{
		Namespace: fmt.Sprintf("%s] - [%s", l.Namespace, ns),
	}
}

func (l *Logger) Info(format string, v ...any) {
	log.Printf("[%s] - [INFO] |> %s", l.Namespace, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(format string, v ...any) {
	log.Printf("[%s] - [ERROR] |> %s", l.Namespace, fmt.Sprintf(format, v...))
}

func (l *Logger) Panic(format string, v ...any) {
	log.Panicf("[%s] - [PANIC] |> %s", l.Namespace, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(format string, v ...any) {
	log.Fatalf("[%s] - [FATAL] |> %s", l.Namespace, fmt.Sprintf(format, v...))
}
