package wasify

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

// asSlogLevel gets 'slog' level based on severity specified by user
func asSlogLevel(severity LogSeverity) slog.Level {
	level, ok := logMap[severity]
	if !ok {
		// default logger is Info
		return slog.LevelInfo
	}

	return level
}

type Logger interface {
	Severity() LogSeverity
	ForSeverity(severity LogSeverity) Logger

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
			Level: asSlogLevel(severity),
		})),
	}
}

type _SlogLogger struct {
	severity LogSeverity
	delegate *slog.Logger
}

func (self *_SlogLogger) Severity() LogSeverity {
	return self.severity
}

func (self *_SlogLogger) ForSeverity(severity LogSeverity) Logger {
	return &_SlogLogger{
		delegate: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: asSlogLevel(severity),
		})),
	}
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
	self.delegate.Log(context.Background(), asSlogLevel(severity), message, arguments...)
}
