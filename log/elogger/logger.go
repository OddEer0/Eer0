package elogger

import (
	"context"
	"io"
)

type (
	Level int

	LogValuer interface {
		LogValue() Value
	}

	Value struct {
		Val interface{}
	}

	Field struct {
		Key   string
		Value Value
	}

	Logger interface {
		WithFields(fields ...Field) Logger
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
