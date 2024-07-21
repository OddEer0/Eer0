package elogger_test

import (
	"context"
	"github.com/OddEer0/Eer0/elogger"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type dumpOut struct {
	res string
}

func (d *dumpOut) Write(p []byte) (n int, err error) {
	d.res = string(p)
	return len(p), nil
}

func TestBase(t *testing.T) {
	t.Run("Should log", func(t *testing.T) {
		dump := &dumpOut{}
		ctx := context.Background()
		log := elogger.New(&elogger.Options{
			Level:  elogger.DebugLevel,
			Output: dump,
		})
		log.Debug(ctx, "hi", elogger.Any("kek", "haha"))
		assert.NotEmpty(t, dump.res)
		log.Info(ctx, "hi")
		assert.NotEmpty(t, dump.res)
		log.Warn(ctx, "hi")
		assert.NotEmpty(t, dump.res)
		log.Error(ctx, "hi")
		assert.NotEmpty(t, dump.res)
	})

	t.Run("Should correct all level write", func(t *testing.T) {
		ctx := context.Background()
		dump := &dumpOut{}
		log := elogger.New(&elogger.Options{
			Level:   elogger.DebugLevel,
			Output:  dump,
			OffTime: true,
		})
		writeStr := "Hello"
		fields := []elogger.Field{{
			Key:   "f1",
			Value: "lol",
		}, {
			Key:   "f2",
			Value: "kek",
		}}
		log.Debug(ctx, writeStr, fields...)
		assert.Equal(t, `{"level":"DEBUG","msg":"Hello","f1":"lol","f2":"kek"}`+"\n", dump.res)
		log.Info(ctx, writeStr, fields...)
		assert.Equal(t, `{"level":"INFO","msg":"Hello","f1":"lol","f2":"kek"}`+"\n", dump.res)
		log.Warn(ctx, writeStr, fields...)
		assert.Equal(t, `{"level":"WARN","msg":"Hello","f1":"lol","f2":"kek"}`+"\n", dump.res)
		log.Error(ctx, writeStr, fields...)
		assert.Equal(t, `{"level":"ERROR","msg":"Hello","f1":"lol","f2":"kek"}`+"\n", dump.res)

		log.Debug(ctx, writeStr)
		assert.Equal(t, `{"level":"DEBUG","msg":"Hello"}`+"\n", dump.res)
		log.Info(ctx, writeStr)
		assert.Equal(t, `{"level":"INFO","msg":"Hello"}`+"\n", dump.res)
		log.Warn(ctx, writeStr)
		assert.Equal(t, `{"level":"WARN","msg":"Hello"}`+"\n", dump.res)
		log.Error(ctx, writeStr)
		assert.Equal(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dump.res)
	})

	t.Run("Should correct level output", func(t *testing.T) {
		ctx := context.Background()
		dump := &dumpOut{}
		dumpError := &dumpOut{}
		dumpWarn := &dumpOut{}
		log := elogger.New(&elogger.Options{
			Level:  elogger.DebugLevel,
			Output: dump,
			LevelOutput: map[elogger.Level]io.Writer{
				elogger.ErrorLevel: dumpError,
				elogger.WarnLevel:  dumpWarn,
			},
			OffTime: true,
		})
		writeStr := "Hello"
		log.Debug(ctx, writeStr)
		assert.Equal(t, `{"level":"DEBUG","msg":"Hello"}`+"\n", dump.res)
		log.Info(ctx, writeStr)
		assert.Equal(t, `{"level":"INFO","msg":"Hello"}`+"\n", dump.res)
		log.Warn(ctx, writeStr)
		assert.NotEqual(t, `{"level":"WARN","msg":"Hello"}`+"\n", dump.res)
		assert.Equal(t, `{"level":"WARN","msg":"Hello"}`+"\n", dumpWarn.res)
		assert.NotEqual(t, `{"level":"WARN","msg":"Hello"}`+"\n", dumpError.res)
		log.Error(ctx, writeStr)
		assert.NotEqual(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dump.res)
		assert.NotEqual(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dumpWarn.res)
		assert.Equal(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dumpError.res)
	})

	t.Run("Should correct log level", func(t *testing.T) {
		dump := &dumpOut{}
		opt := &elogger.Options{
			Level:   elogger.InfoLevel,
			Output:  dump,
			OffTime: true,
		}
		ctx := context.Background()
		log := elogger.New(opt)
		writeStr := "Hello"

		log.Debug(ctx, writeStr)
		assert.Equal(t, "", dump.res)
		log.Info(ctx, writeStr)
		assert.Equal(t, `{"level":"INFO","msg":"Hello"}`+"\n", dump.res)
		log.Warn(ctx, writeStr)
		assert.Equal(t, `{"level":"WARN","msg":"Hello"}`+"\n", dump.res)
		log.Error(ctx, writeStr)
		assert.Equal(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dump.res)
		dump.res = ""
		opt.Level = elogger.WarnLevel
		log.Debug(ctx, writeStr)
		assert.Equal(t, "", dump.res)
		log.Info(ctx, writeStr)
		assert.Equal(t, "", dump.res)
		log.Warn(ctx, writeStr)
		assert.Equal(t, `{"level":"WARN","msg":"Hello"}`+"\n", dump.res)
		log.Error(ctx, writeStr)
		assert.Equal(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dump.res)
		dump.res = ""
		opt.Level = elogger.ErrorLevel
		assert.Equal(t, "", dump.res)
		log.Info(ctx, writeStr)
		assert.Equal(t, "", dump.res)
		log.Warn(ctx, writeStr)
		assert.Equal(t, "", dump.res)
		log.Error(ctx, writeStr)
		assert.Equal(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dump.res)
	})
}
