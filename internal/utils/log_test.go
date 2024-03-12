package utils

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogLevel(t *testing.T) {
	newLogger := NewSlogLogger(LogDebug)
	assert.NotNil(t, newLogger)

	tests := []struct {
		severity LogSeverity
		expected slog.Level
	}{
		{LogDebug, slog.LevelDebug},
		{LogInfo, slog.LevelInfo},
		{LogWarning, slog.LevelWarn},
		{LogError, slog.LevelError},
		{LogSeverity(255), slog.LevelInfo}, // Unexpected severity
	}

	for _, test := range tests {
		got := GetlogLevel(test.severity)
		if got != test.expected {
			t.Errorf("for severity %d, expected %d but got %d", test.severity, test.expected, got)
		}
	}
}
