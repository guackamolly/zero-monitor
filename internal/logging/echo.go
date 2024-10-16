package logging

import "github.com/labstack/echo/v4"

type echoLogger struct {
	logger echo.Logger
}

func (l echoLogger) Info(fmt string, s ...any) {
	l.logger.Infof(fmt, s...)
}

func (l echoLogger) Warning(fmt string, s ...any) {
	l.logger.Warnf(fmt, s...)
}

func (l echoLogger) Error(fmt string, s ...any) {
	l.logger.Errorf(fmt, s...)
}

func (l echoLogger) Fatal(fmt string, s ...any) {
	l.logger.Fatalf(fmt, s...)
}

func NewEchoLogger(
	logger echo.Logger,
) Logger {
	return echoLogger{
		logger: logger,
	}
}
