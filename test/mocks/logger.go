package mocks

import (
	"context"
	"log/slog"
)

// handler implements slog.Handler interface.
type handler struct{}

// Ensure that handler implements slog.Handler interface.
var _ slog.Handler = handler{}

// Ensure that handler implements slog.Handler interface.
var _ slog.Handler = handler{}

// NewNullLogger returns a logger that discards all log messages.
func NewNullLogger() slog.Logger {
	return *slog.New(handler{})
}

// Enabled always returns false, meaning no log levels are enabled.
func (h handler) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}

// Handle discards the log record, effectively doing nothing.
func (h handler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

// WithAttrs returns a new handler, ignoring the provided attributes.
func (h handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return handler{}
}

// WithGroup returns a new handler, ignoring the provided group name.
func (h handler) WithGroup(name string) slog.Handler {
	return handler{}
}
