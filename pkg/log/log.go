package log

import (
	"fmt"
	"log/slog"
	"os"
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

func Init() {
	opts := slog.HandlerOptions{
		// Level: slog.LevelDebug, // Optional: Set your desired log level
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Check if the attribute key is the time key
			if a.Key == slog.TimeKey {
				// Return an empty Attr, effectively removing it
				return slog.Attr{}
			}
			// Keep all other attributes
			return a
		},
	}

	// Create a TextHandler with the custom options
	logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))

	// Use the logger (or set it as the default)
	slog.SetDefault(logger)
}
