package logger

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger initializes and returns a zap logger.
// The logger is configured based on the provided log level.
func InitLogger(logLevel string, proxyID string) logr.Logger {
	// Parse the log level string to zapcore.Level
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(logLevel)); err != nil {
		// If parsing fails, default to Info level
		zapLevel = zapcore.InfoLevel
	}

	// Create a zap production configuration
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zapLevel)
	zapConfig.InitialFields = map[string]interface{}{
		"proxyID": proxyID,
	}

	// Build the zap logger
	zapLogger, err := zapConfig.Build()
	if err != nil {
		// If building the logger fails, fallback to a no-op logger
		zapLogger = zap.NewNop()
	}

	// Wrap the zap logger with logr interface using zapr
	return zapr.NewLogger(zapLogger)
}
