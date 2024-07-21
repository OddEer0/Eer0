package elogger

import (
	"context"
	"io"
	"log/slog"
)

type (
	Level int

	Field struct {
		Key   string
		Value interface{}
	}

	Logger interface {
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

	logger struct {
		log *slog.Logger
	}
)

func New(opt *Options) Logger {
	h := newHandler(opt)
	l := slog.New(h)
	return &logger{
		log: l,
	}
}
