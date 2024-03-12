package utils

import (
	"context"
	"log/slog"
	"os"
)

type LogSeverity uint8

const (
	LogDebug LogSeverity = iota + 1
	LogInfo
	LogWarning
	LogError
)

var logMap = map[LogSeverity]slog.Level{
	LogDebug:   slog.LevelDebug,
	LogInfo:    slog.LevelInfo,
	LogWarning: slog.LevelWarn,
	LogError:   slog.LevelError,
}

type Logger interface {
	Info(message string, arguments ...any)
	Warn(message string, arguments ...any)
	Error(message string, arguments ...any)
	Debug(message string, arguments ...any)
	Log(severity LogSeverity, message string, arguments ...any)
}

// NewSlogLogger returns new slog ref
func NewSlogLogger(severity LogSeverity) Logger {
	return &_SlogLogger{
		delegate: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: GetlogLevel(severity),
		})),
	}
}

type _SlogLogger struct {
	delegate *slog.Logger
}

func (self *_SlogLogger) Info(message string, arguments ...any) {
	self.Log(LogInfo, message, arguments...)
}

func (self *_SlogLogger) Error(message string, arguments ...any) {
	self.Log(LogError, message, arguments...)
}

func (self *_SlogLogger) Debug(message string, arguments ...any) {
	self.Log(LogDebug, message, arguments...)
}

func (self *_SlogLogger) Warn(message string, arguments ...any) {
	self.Log(LogWarning, message, arguments...)
}

func (self *_SlogLogger) Log(severity LogSeverity, message string, arguments ...any) {
	self.delegate.Log(context.Background(), GetlogLevel(severity), message, arguments...)
}

// GetlogLevel gets 'slog' level based on severity specified by user
func GetlogLevel(s LogSeverity) slog.Level {

	val, ok := logMap[s]
	if !ok {
		// default logger is Info
		return logMap[2]
	}

	return val
}
