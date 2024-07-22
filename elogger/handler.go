package elogger

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"
)

type HandlerOptions struct {
	Level
	LevelOutput map[Level]io.Writer
	Output      io.Writer
	OffTime     bool
	Indent      bool
}

type handler struct {
	opt    *HandlerOptions
	fields []Field
}

func (h handler) Enabled(_ context.Context, level Level) bool {
	return level >= h.opt.Level
}

func (h handler) writeKV(str *strings.Builder, key, value string) {
	if str.Len() > 0 {
		str.WriteString(`,`)
	} else {
		str.WriteString(`{`)
	}
	if h.opt.Indent {
		str.WriteString("\n  ")
	}
	str.WriteString(`"`)
	str.WriteString(key)
	if h.opt.Indent {
		str.WriteString(`": "`)
	} else {
		str.WriteString(`":"`)
	}
	str.WriteString(value)
	str.WriteString(`"`)
}

func (h handler) print(b []byte, level Level) error {
	if o, ok := h.opt.LevelOutput[level]; ok {
		_, err := o.Write(b)
		if err != nil {
			return err
		}
		return nil
	}
	_, err := h.opt.Output.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (h handler) fieldByte(fields []Field) ([]byte, error) {
	res := make(map[string]interface{}, len(h.fields)+len(fields))
	if len(h.fields) > 0 {
		for _, f := range h.fields {
			if _, ok := res[f.Key]; !ok {
				res[f.Key] = f.Value
			}
		}
	}
	for _, f := range fields {
		if _, ok := res[f.Key]; !ok {
			res[f.Key] = f.Value
		}
	}

	if h.opt.Indent {
		return json.MarshalIndent(res, "", "  ")
	}
	return json.Marshal(res)
}

func (h handler) concat(str *strings.Builder, fields []byte) []byte {
	if len(fields) <= 7 {
		return []byte(str.String())
	}
	res := make([]byte, 0, len(fields)+str.Len())
	res = append(res, []byte(str.String())...)
	if h.opt.Indent {
		res = res[:len(res)-1]
	}
	res[len(res)-1] = ','
	res = append(res, fields[1:]...)
	return res
}

func (h handler) Handle(ctx context.Context, r Record) error {
	res := &strings.Builder{}
	if !h.opt.OffTime {
		timeStr := r.Time.Format(time.DateTime)
		h.writeKV(res, "time", timeStr)
	}
	h.writeKV(res, "level", r.Level.String())
	h.writeKV(res, "msg", r.Message)
	if h.opt.Indent {
		res.WriteString("\n")
	}
	res.WriteString("}")

	fieldByte, err := h.fieldByte(r.fields)
	if err != nil {
		return err
	}
	writeByte := h.concat(res, fieldByte)
	writeByte = append(writeByte, '\n')

	err = h.print(writeByte, r.Level)
	if err != nil {
		return err
	}
	return nil
}

func (h handler) WithFields(attrs []Field) Handler {
	f := make([]Field, 0, len(h.fields)+len(attrs))
	if len(h.fields) > 0 {
		f = append(f, h.fields...)
	}
	f = append(f, attrs...)
	return &handler{
		opt:    h.opt,
		fields: f,
	}
}

func NewHandler(opt *HandlerOptions) Handler {
	return &handler{
		opt: opt,
	}
}
