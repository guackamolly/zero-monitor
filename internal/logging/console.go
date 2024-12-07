package logging

import (
	"fmt"
	console "log"

	"github.com/labstack/gommon/color"
)

var debugPrefix = color.Green("(debug): ")
var infoPrefix = color.Blue("(info): ")
var warningPrefix = color.Yellow("(warning): ")
var errorPrefix = color.Red("(error): ")
var fatalPrefix = color.RedBg("(fatal): ")

type consoleLogger struct{}

func (l consoleLogger) Info(fmt string, s ...any) {
	console.Println(infoPrefix + l.format(fmt, s...))
}

func (l consoleLogger) Warning(fmt string, s ...any) {
	console.Println(warningPrefix + l.format(fmt, s...))
}

func (l consoleLogger) Error(fmt string, s ...any) {
	console.Println(errorPrefix + l.format(fmt, s...))
}

func (l consoleLogger) Fatal(fmt string, s ...any) {
	f := l.format(fmt, s...)
	console.Fatalln(fatalPrefix + f)
}

func (l consoleLogger) Debug(fmt string, s ...any) {
	console.Println(debugPrefix + l.format(fmt, s...))
}

func (l consoleLogger) format(fmts string, s ...any) string {
	return fmt.Sprintf(fmts, s...)
}

func NewConsoleLogger() Logger {
	return consoleLogger{}
}
