package slog

import (
	"context"
	"log/slog"

	"github.com/pasataleo/go-logging/pkg/logging"
)

// Handler is an slog.Handler that forwards records to a logging.Logger.
type Handler struct {
	logger logging.Logger
	group  string
	attrs  []slog.Attr
}

// NewHandler creates an slog.Handler backed by the given logging.Logger.
func NewHandler(logger logging.Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return h.logger.ShouldLog(mapLevel(level))
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	l := h.logger

	// Apply pre-set attrs.
	for _, attr := range h.attrs {
		l = l.WithField(h.prefixKey(attr.Key), attr.Value.Any())
	}

	// Apply record attrs.
	record.Attrs(func(attr slog.Attr) bool {
		l = l.WithField(h.prefixKey(attr.Key), attr.Value.Any())
		return true
	})

	level := mapLevel(record.Level)
	msg := record.Message

	switch {
	case level >= logging.Panic:
		l.Panic(msg)
	case level >= logging.Fatal:
		l.Fatal(msg)
	case level >= logging.Error:
		l.Error(msg)
	case level >= logging.Warn:
		l.Warn(msg)
	case level >= logging.Info:
		l.Info(msg)
	case level >= logging.Debug:
		l.Debug(msg)
	default:
		l.Trace(msg)
	}

	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	combined := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(combined, h.attrs)
	copy(combined[len(h.attrs):], attrs)
	return &Handler{
		logger: h.logger,
		group:  h.group,
		attrs:  combined,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	g := name
	if h.group != "" {
		g = h.group + "." + name
	}
	return &Handler{
		logger: h.logger,
		group:  g,
		attrs:  h.attrs,
	}
}

func (h *Handler) prefixKey(key string) string {
	if h.group == "" {
		return key
	}
	return h.group + "." + key
}

// mapLevel converts an slog.Level to a logging.Level.
func mapLevel(level slog.Level) logging.Level {
	switch {
	case level >= slog.LevelError:
		return logging.Error
	case level >= slog.LevelWarn:
		return logging.Warn
	case level >= slog.LevelInfo:
		return logging.Info
	case level >= slog.LevelDebug:
		return logging.Debug
	default:
		return logging.Trace
	}
}
