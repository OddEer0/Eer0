package elogger_test

import (
	"context"
	"encoding/json"
	"github.com/OddEer0/Eer0/elogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		h := elogger.NewHandler(&elogger.HandlerOptions{
			Level:  elogger.DebugLevel,
			Output: dump,
		})
		var data interface{}
		log := elogger.New(h)
		log.Debug(ctx, "hi", elogger.Any("kek", "haha"))
		assert.NotEmpty(t, dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Info(ctx, "hi")
		assert.NotEmpty(t, dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Warn(ctx, "hi")
		assert.NotEmpty(t, dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Error(ctx, "hi")
		assert.NotEmpty(t, dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
	})

	t.Run("Should correct all level write", func(t *testing.T) {
		ctx := context.Background()
		dump := &dumpOut{}
		h := elogger.NewHandler(&elogger.HandlerOptions{
			Level:   elogger.DebugLevel,
			Output:  dump,
			OffTime: true,
		})
		var data interface{}
		log := elogger.New(h)
		writeStr := "Hello"
		fields := []elogger.Field{elogger.Any(
			"f1",
			"lol",
		), elogger.Any(
			"f2",
			"kek",
		)}
		log.Debug(ctx, writeStr, fields...)
		assert.Equal(t, `{"level":"DEBUG","msg":"Hello","f1":"lol","f2":"kek"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Info(ctx, writeStr, fields...)
		assert.Equal(t, `{"level":"INFO","msg":"Hello","f1":"lol","f2":"kek"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Warn(ctx, writeStr, fields...)
		assert.Equal(t, `{"level":"WARN","msg":"Hello","f1":"lol","f2":"kek"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Error(ctx, writeStr, fields...)
		assert.Equal(t, `{"level":"ERROR","msg":"Hello","f1":"lol","f2":"kek"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))

		log.Debug(ctx, writeStr)
		assert.Equal(t, `{"level":"DEBUG","msg":"Hello"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Info(ctx, writeStr)
		assert.Equal(t, `{"level":"INFO","msg":"Hello"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Warn(ctx, writeStr)
		assert.Equal(t, `{"level":"WARN","msg":"Hello"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Error(ctx, writeStr)
		assert.Equal(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
	})

	t.Run("Should correct level output", func(t *testing.T) {
		ctx := context.Background()
		dump := &dumpOut{}
		dumpError := &dumpOut{}
		dumpWarn := &dumpOut{}
		h := elogger.NewHandler(&elogger.HandlerOptions{
			Level:   elogger.DebugLevel,
			Output:  dump,
			OffTime: true,
			LevelOutput: map[elogger.Level]io.Writer{
				elogger.ErrorLevel: dumpError,
				elogger.WarnLevel:  dumpWarn,
			},
		})
		var data interface{}
		log := elogger.New(h)
		writeStr := "Hello"
		log.Debug(ctx, writeStr)
		assert.Equal(t, `{"level":"DEBUG","msg":"Hello"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Info(ctx, writeStr)
		assert.Equal(t, `{"level":"INFO","msg":"Hello"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Warn(ctx, writeStr)
		assert.NotEqual(t, `{"level":"WARN","msg":"Hello"}`+"\n", dump.res)
		assert.Equal(t, `{"level":"WARN","msg":"Hello"}`+"\n", dumpWarn.res)
		require.NoError(t, json.Unmarshal([]byte(dumpWarn.res), &data))
		assert.NotEqual(t, `{"level":"WARN","msg":"Hello"}`+"\n", dumpError.res)
		log.Error(ctx, writeStr)
		assert.NotEqual(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dump.res)
		assert.NotEqual(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dumpWarn.res)
		assert.Equal(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dumpError.res)
		require.NoError(t, json.Unmarshal([]byte(dumpError.res), &data))
	})

	t.Run("Should correct log level", func(t *testing.T) {
		dump := &dumpOut{}
		opt := &elogger.HandlerOptions{
			Level:   elogger.InfoLevel,
			Output:  dump,
			OffTime: true,
		}
		h := elogger.NewHandler(opt)
		ctx := context.Background()
		log := elogger.New(h)
		writeStr := "Hello"
		var data interface{}

		log.Debug(ctx, writeStr)
		assert.Equal(t, "", dump.res)
		log.Info(ctx, writeStr)
		assert.Equal(t, `{"level":"INFO","msg":"Hello"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Warn(ctx, writeStr)
		assert.Equal(t, `{"level":"WARN","msg":"Hello"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Error(ctx, writeStr)
		assert.Equal(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		dump.res = ""
		opt.Level = elogger.WarnLevel
		log.Debug(ctx, writeStr)
		assert.Equal(t, "", dump.res)
		log.Info(ctx, writeStr)
		assert.Equal(t, "", dump.res)
		log.Warn(ctx, writeStr)
		assert.Equal(t, `{"level":"WARN","msg":"Hello"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Error(ctx, writeStr)
		assert.Equal(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		dump.res = ""
		opt.Level = elogger.ErrorLevel
		assert.Equal(t, "", dump.res)
		log.Info(ctx, writeStr)
		assert.Equal(t, "", dump.res)
		log.Warn(ctx, writeStr)
		assert.Equal(t, "", dump.res)
		log.Error(ctx, writeStr)
		assert.Equal(t, `{"level":"ERROR","msg":"Hello"}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
	})

	t.Run("Should correct indent", func(t *testing.T) {
		dump := &dumpOut{}
		opt := &elogger.HandlerOptions{
			Level:   elogger.DebugLevel,
			Output:  dump,
			OffTime: true,
			Indent:  true,
		}
		var data interface{}
		h := elogger.NewHandler(opt)
		ctx := context.Background()
		log := elogger.New(h)
		writeStr := "Hello"

		log.Debug(ctx, writeStr)
		assert.Equal(t, `{
  "level": "DEBUG",
  "msg": "Hello"
}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
		log.Info(ctx, writeStr, elogger.Any("lol", "kek"), elogger.Any("kek", "lol"))
		assert.Equal(t, `{
  "level": "INFO",
  "msg": "Hello",
  "kek": "lol",
  "lol": "kek"
}`+"\n", dump.res)
		require.NoError(t, json.Unmarshal([]byte(dump.res), &data))
	})
}
