package elogger

import (
	"context"
	"time"
)

func (l *logger) Log(ctx context.Context, level Level, msg string, fields ...Field) {
	if !l.handler.Enabled(ctx, level) {
		return
	}

	err := l.handler.Handle(ctx, Record{
		Level:   level,
		Message: msg,
		Time:    time.Now(),
		fields:  fields,
	})

	if err != nil {
		panic(err)
	}
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.Log(ctx, DebugLevel, msg, fields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...Field) {
	l.Log(ctx, InfoLevel, msg, fields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.Log(ctx, WarnLevel, msg, fields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...Field) {
	l.Log(ctx, ErrorLevel, msg, fields...)
}
