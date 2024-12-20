package logging

var loggers = []Logger{}
var includeDebugLogs = true

func LogInfo(fmt string, s ...any) {
	for _, l := range loggers {
		log(l.Info, fmt, s...)
	}
}

func LogWarning(fmt string, s ...any) {
	for _, l := range loggers {
		log(l.Warning, fmt, s...)
	}
}

func LogError(fmt string, s ...any) {
	for _, l := range loggers {
		log(l.Error, fmt, s...)
	}
}

func LogFatal(fmt string, s ...any) {
	for _, l := range loggers {
		log(l.Fatal, fmt, s...)
	}
}

func LogDebug(fmt string, s ...any) {
	if !includeDebugLogs {
		return
	}

	for _, l := range loggers {
		log(l.Debug, fmt, s...)
	}
}

func log(cb func(string, ...any), fmt string, s ...any) {
	cb(fmt, s...)
}

type Logger interface {
	Info(fmt string, s ...any)
	Warning(fmt string, s ...any)
	Error(fmt string, s ...any)
	Fatal(fmt string, s ...any)
	Debug(fmt string, s ...any)
}

func AddLogger(logger Logger) {
	loggers = append(loggers, logger)
}

func DisableDebugLogs() {
	includeDebugLogs = false
}
