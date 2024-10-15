package elogger_test

import (
	"context"
	"github.com/OddEer0/Eer0/log/elogger"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	User struct {
		Id       string `json:"id"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	UserL struct {
		Id       string `json:"id"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}
)

func (u UserL) LogValue() elogger.Value {
	return elogger.GroupValue(
		elogger.Any("id", u.Id),
		elogger.Any("login", u.Login),
	)
}

type AppHandler struct {
	elogger.Handler
}

func (h AppHandler) Handle(ctx context.Context, r elogger.Record) error {
	r.AddFields(elogger.Any("kek", "lol"))
	r.For(func(field elogger.Field) bool {
		if field.Key == "wait" {
			r.AddFields(elogger.Any("sec", "42"))
			return false
		}
		return true
	})
	r.AddFields(elogger.Any("count", r.LenFields()))
	return h.Handler.Handle(ctx, r)
}

func TestCustom(t *testing.T) {
	user := User{Id: "1", Login: "kek", Password: "root"}
	userL := UserL{Id: "1", Login: "kek", Password: "root"}

	t.Run("Should LogValuer interface", func(t *testing.T) {
		dump := &dumpOut{}
		ctx := context.Background()
		h := elogger.NewHandler(&elogger.HandlerOptions{
			Level:   elogger.DebugLevel,
			Output:  dump,
			OffTime: true,
		})
		log := elogger.New(h)
		writeText := "Hello"
		log.Info(ctx, writeText, elogger.Any("user", user))
		assert.Equal(t, `{"level":"INFO","msg":"Hello","user":{"id":"1","login":"kek","password":"root"}}`+"\n", dump.res)
		log.Info(ctx, writeText, elogger.Any("user", userL))
		assert.Equal(t, `{"level":"INFO","msg":"Hello","user":{"id":"1","login":"kek"}}`+"\n", dump.res)
	})

	t.Run("Should custom handler", func(t *testing.T) {
		dump := &dumpOut{}
		ctx := context.Background()
		h := AppHandler{
			elogger.NewHandler(&elogger.HandlerOptions{
				Level:   elogger.DebugLevel,
				Output:  dump,
				OffTime: true,
			}),
		}
		log := elogger.New(h)
		writeText := "Hello"
		log.Info(ctx, writeText)
		assert.Equal(t, `{"level":"INFO","msg":"Hello","count":1,"kek":"lol"}`+"\n", dump.res)
		log.Info(ctx, writeText, elogger.Any("wait", "second"))
		assert.Equal(t, `{"level":"INFO","msg":"Hello","count":3,"kek":"lol","sec":"42","wait":"second"}`+"\n", dump.res)
	})

	t.Run("Should logger with field", func(t *testing.T) {
		dump := &dumpOut{}
		ctx := context.Background()
		h := elogger.NewHandler(&elogger.HandlerOptions{
			Level:   elogger.DebugLevel,
			Output:  dump,
			OffTime: true,
		})
		log := elogger.New(h)
		log = log.WithFields(elogger.Group("in", elogger.Any("1", 1), elogger.Any("2", 2)))
		writeText := "Hello"
		log.Info(ctx, writeText)
		assert.Equal(t, `{"level":"INFO","msg":"Hello","in":{"1":1,"2":2}}`+"\n", dump.res)

		log = log.WithFields(elogger.Any("hah", "haha"))
		log.Info(ctx, writeText)
		assert.Equal(t, `{"level":"INFO","msg":"Hello","hah":"haha","in":{"1":1,"2":2}}`+"\n", dump.res)
	})
}
