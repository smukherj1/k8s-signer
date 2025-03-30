package log

import (
	"fmt"
	"log/slog"
)

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Debugf(msg string, args ...any) {
	Debug(fmt.Sprintf(msg, args...))
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Infof(msg string, args ...any) {
	Info(fmt.Sprintf(msg, args...))
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Warnf(msg string, args ...any) {
	Warn(fmt.Sprintf(msg, args...))
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Errorf(msg string, args ...any) {
	Error(fmt.Sprintf(msg, args...))
}
