package logger

import (
	"fmt"
	"log/slog"
)

// SlogAdapter adapts *slog.Logger to satisfy runtime.Logger interface (Printf method).
// Pattern: Adapter - bridge between runtime.Logger and slog.Logger
type SlogAdapter struct {
	logger *slog.Logger
}

// NewSlog creates a new slog adapter.
func NewSlog(logger *slog.Logger) *SlogAdapter {
	if logger == nil {
		logger = slog.Default()
	}
	return &SlogAdapter{logger: logger}
}

// Printf implements runtime.Logger interface.
// Converts Printf-style logging to structured slog.InfoContext.
func (a *SlogAdapter) Printf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	a.logger.Info(message)
}

// Info logs an info-level message with context.
func (a *SlogAdapter) Info(msg string) {
	a.logger.Info(msg)
}

// Infof logs a formatted info-level message.
func (a *SlogAdapter) Infof(format string, args ...interface{}) {
	a.logger.Info(fmt.Sprintf(format, args...))
}

// Error logs an error-level message.
func (a *SlogAdapter) Error(msg string) {
	a.logger.Error(msg)
}

// Errorf logs a formatted error-level message.
func (a *SlogAdapter) Errorf(format string, args ...interface{}) {
	a.logger.Error(fmt.Sprintf(format, args...))
}

// Debug logs a debug-level message.
func (a *SlogAdapter) Debug(msg string) {
	a.logger.Debug(msg)
}

// Debugf logs a formatted debug-level message.
func (a *SlogAdapter) Debugf(format string, args ...interface{}) {
	a.logger.Debug(fmt.Sprintf(format, args...))
}

// WithAttrs returns a new SlogAdapter with the given attributes.
// Useful for structured logging context.
func (a *SlogAdapter) WithAttrs(attrs ...slog.Attr) *SlogAdapter {
	// Convert attrs to any slice for With method
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	return &SlogAdapter{logger: a.logger.With(args...)}
}

// Underlying returns the underlying *slog.Logger for advanced use cases.
func (a *SlogAdapter) Underlying() *slog.Logger {
	return a.logger
}
