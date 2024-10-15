package elogger

import (
	"time"
)

type Record struct {
	Time time.Time
	Level
	Message string
	fields  []Field
}

func (r *Record) AddFields(fields ...Field) {
	r.fields = append(r.fields, fields...)
}

func (r *Record) LenFields() int {
	return len(r.fields)
}

func (r *Record) For(fn func(Field) bool) {
	for _, field := range r.fields {
		if !fn(field) {
			break
		}
	}
}
