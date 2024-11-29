package logging

import (
	"fmt"
	console "log"
)

type consoleLogger struct{}

func (l consoleLogger) Info(fmt string, s ...any) {
	console.Println("(info): " + l.format(fmt, s...))
}

func (l consoleLogger) Warning(fmt string, s ...any) {
	console.Println("(warn): " + l.format(fmt, s...))
}

func (l consoleLogger) Error(fmt string, s ...any) {
	console.Println("(error): " + l.format(fmt, s...))
}

func (l consoleLogger) Fatal(fmt string, s ...any) {
	f := l.format(fmt, s...)
	console.Fatalln("(fatal): " + f)
}

func (l consoleLogger) Debug(fmt string, s ...any) {
	console.Println("(debug): " + l.format(fmt, s...))
}

func (l consoleLogger) format(fmts string, s ...any) string {
	return fmt.Sprintf(fmts, s...)
}

func NewConsoleLogger() Logger {
	return consoleLogger{}
}
