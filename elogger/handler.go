package elogger

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"
)

type handler struct {
	opt   *Options
	attrs []slog.Attr
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= mapLevelToSlogLevel(h.opt.Level)
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})
	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	timeStr := r.Time.Format(time.DateTime)
	str := strings.Builder{}
	if !h.opt.OffTime {
		str.WriteString(`{"time":"`)
		str.WriteString(timeStr)
		str.WriteString(`","level":"`)
		str.WriteString(r.Level.String())
	} else {
		str.WriteString(`{"level":"`)
		str.WriteString(r.Level.String())
	}
	str.WriteString(`","msg":"`)
	str.WriteString(r.Message)
	str.WriteString(`"`)

	var b []byte
	var err error
	if len(fields) > 0 {
		str.WriteString(",")
		b, err = json.Marshal(fields)
		//b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
		b = b[1:]
	} else {
		str.WriteString("}")
	}
	out := h.opt.Output

	if o, ok := h.opt.LevelOutput[mapSlogLevelToLevel(r.Level)]; ok {
		out = o
	}
	res := make([]byte, 0, len(b)+str.Len()+10)
	res = append(res, []byte(str.String())...)
	res = append(res, b...)
	res = append(res, '\n')

	_, err = out.Write(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newH := &handler{
		attrs: attrs,
		opt:   h.opt,
	}
	if h.attrs != nil {
		newH.attrs = append(newH.attrs, h.attrs...)
	}
	return newH
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		opt:   h.opt,
		attrs: h.attrs,
	}
}

func newHandler(opt *Options) slog.Handler {
	return &handler{
		opt: opt,
	}
}
