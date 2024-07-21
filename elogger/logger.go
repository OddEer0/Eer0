package elogger

import (
	"context"
	"io"
)

type (
	Level int

	Field struct {
		Key   string
		Value interface{}
	}

	Logger interface {
		Log(ctx context.Context, level Level, format string, fields ...Field)
		Debug(ctx context.Context, msg string, fields ...Field)
		Info(ctx context.Context, msg string, fields ...Field)
		Warn(ctx context.Context, msg string, fields ...Field)
		Error(ctx context.Context, msg string, fields ...Field)
	}

	Options struct {
		OffTime     bool
		Output      io.Writer
		LevelOutput map[Level]io.Writer
		Level
	}

	Handler interface {
		Enabled(context.Context, Level) bool
		Handle(context.Context, Record) error
		WithFields(attrs []Field) Handler
	}

	logger struct {
		handler Handler
	}
)

func New(h Handler) Logger {
	return &logger{
		handler: h,
	}
}
