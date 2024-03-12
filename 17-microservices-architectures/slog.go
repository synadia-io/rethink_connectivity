package main

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
)

type natsLogWriter struct {
	subject string
	nc      *nats.Conn
}

// Write implements io.Writer.
func (n *natsLogWriter) Write(p []byte) (int, error) {
	err := n.nc.Publish(n.subject, p)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

type natsSlogHandler struct {
	handler slog.Handler
	level   slog.Level
}

// Enabled implements slog.Handler.
func (n *natsSlogHandler) Enabled(c context.Context, l slog.Level) bool {
	return n.handler.Enabled(c, l)
}

// Handle implements slog.Handler.
func (n *natsSlogHandler) Handle(c context.Context, record slog.Record) error {
	if record.Level < n.level {
		return nil
	}
	return n.handler.Handle(c, record)
}

// WithAttrs implements slog.Handler.
func (n *natsSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &natsSlogHandler{n.handler.WithAttrs(attrs), n.level}
}

// WithGroup implements slog.Handler.
func (n *natsSlogHandler) WithGroup(name string) slog.Handler {
	return &natsSlogHandler{n.handler.WithGroup(name), n.level}
}
