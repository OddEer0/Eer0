package elogger

import (
	"context"
	"log/slog"
)

func (l *logger) Debug(ctx context.Context, msg string, fields ...Field) {
	attrs := make([]any, 0, len(fields))
	for _, field := range fields {
		attrs = append(attrs, slog.Any(field.Key, field.Value))
	}
	l.log.DebugContext(ctx, msg, attrs...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...Field) {
	attrs := make([]any, 0, len(fields))
	for _, field := range fields {
		attrs = append(attrs, slog.Any(field.Key, field.Value))
	}
	l.log.InfoContext(ctx, msg, attrs...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...Field) {
	attrs := make([]any, 0, len(fields))
	for _, field := range fields {
		attrs = append(attrs, slog.Any(field.Key, field.Value))
	}
	l.log.WarnContext(ctx, msg, attrs...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...Field) {
	attrs := make([]any, 0, len(fields))
	for _, field := range fields {
		attrs = append(attrs, slog.Any(field.Key, field.Value))
	}
	l.log.ErrorContext(ctx, msg, attrs...)
}
