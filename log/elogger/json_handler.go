package elogger

import (
	"context"
	"io"
	"log/slog"
)

type jsonHandler struct {
	slog.Handler
	fields []Field
}

type JsonHandlerOption struct {
	AddSource bool
	Level
	Output io.Writer
}

func (j *jsonHandler) Enabled(ctx context.Context, level Level) bool {
	return j.Handler.Enabled(ctx, mapLevelToSlogLevel(level))
}

func (j *jsonHandler) Handle(ctx context.Context, r Record) error {
	var pc uintptr
	record := slog.NewRecord(r.Time, mapLevelToSlogLevel(r.Level), r.Message, pc)
	attrs := make([]slog.Attr, 0, len(r.fields))
	for _, field := range r.fields {
		attrs = append(attrs, slog.Any(field.Key, field.Value))
	}
	record.AddAttrs(attrs...)
	return j.Handler.Handle(ctx, record)
}

func (j *jsonHandler) WithFields(fields []Field) Handler {
	f := make([]Field, len(fields)+len(j.fields))
	copy(f, fields)
	if j.fields != nil {
		f = append(f, j.fields...)
	}
	return &jsonHandler{
		Handler: j.Handler,
		fields:  f,
	}
}

func NewJsonHandler(opt *JsonHandlerOption) Handler {
	return &jsonHandler{
		Handler: slog.NewJSONHandler(opt.Output, &slog.HandlerOptions{
			Level:     mapLevelToSlogLevel(opt.Level),
			AddSource: opt.AddSource,
		}),
	}
}
