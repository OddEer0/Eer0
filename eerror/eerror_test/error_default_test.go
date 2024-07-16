package eerror_test

import (
	"errors"
	"github.com/OddEer0/Eer0/eerror"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func TestDefaultError(t *testing.T) {
	t.Run("Should correct change default error status code", func(t *testing.T) {
		defer eerror.DefaultAll()
		code := eerror.Status(errors.New("another"))
		assert.Equal(t, eerror.NotCode, code)
		err := eerror.New("orig err").Err()
		assert.Equal(t, eerror.DefaultErrorCode, eerror.Status(err))

		eerror.DefaultCode(404)
		err = eerror.New("orig err").Err()
		assert.Equal(t, eerror.ErrorCode(404), eerror.Status(err))
	})

	t.Run("Should correct change default ErrorMethod", func(t *testing.T) {
		defer eerror.DefaultAll()
		err := eerror.New("orig err").Code(100).Err()
		assert.Equal(t, "orig err", err.Error())

		eerror.DefaultErrorMethod(func(e *eerror.Error) string {
			b := strings.Builder{}
			errStatus := strconv.Itoa(int(e.StatusCode))
			b.WriteString("{status: ")
			b.WriteString(errStatus)
			b.WriteString("} ")
			b.WriteString(e.Message)
			return b.String()
		})
		assert.Equal(t, "{status: 100} orig err", err.Error())
	})

	t.Run("Should correct change default Stack capacity", func(t *testing.T) {
		defer eerror.DefaultAll()
		err := eerror.New("orig err").Stack("1").Err()
		var bErr *eerror.Error
		if errors.As(err, &bErr) {
			assert.Equal(t, 8, bErr.StackTrace.Cap())
		}

		eerror.DefaultStackDepth(12)
		err = eerror.New("orig err").Stack("1").Err()
		if errors.As(err, &bErr) {
			assert.Equal(t, 12, bErr.StackTrace.Cap())
		}
	})

	t.Run("Should correct change info cap", func(t *testing.T) {
		defer eerror.DefaultAll()
		err := eerror.New("orig err").Err()
		bErr, ok := err.(*eerror.Error)
		assert.True(t, ok)
		assert.Nil(t, bErr.Info)

		eerror.DefaultInfoSize(32)

		err = eerror.Err(err).
			Field("first", "val").
			Err()
		bErr, ok = err.(*eerror.Error)
		assert.True(t, ok)
		assert.NotNil(t, bErr.Info)
	})

	t.Run("Should correct change cause wrap separator", func(t *testing.T) {
		defer eerror.DefaultAll()
		err := eerror.New("orig err").
			Wrap("[first]").
			Wrap("[second]").
			Err()

		causeStr := "[second]" + eerror.DefaultWrapSeparator + "[first]" + eerror.DefaultWrapSeparator + "orig err"
		assert.Equal(t, causeStr, eerror.Cause(err).Error())

		eerror.DefaultCauseWrapSeparator("=> ")
		err = eerror.New("orig err").
			Wrap("[first]").
			Wrap("[second]").
			Err()

		assert.Equal(t, "[second]=> [first]=> orig err", eerror.Cause(err).Error())
	})

	t.Run("Should correct cause cap", func(t *testing.T) {
		defer eerror.DefaultAll()
		err := eerror.New("orig err").
			Wrap("[first]").
			Wrap("[second]").
			Err()
		wErr, ok := eerror.Cause(err).(*eerror.WrapError)
		assert.True(t, ok)
		assert.Equal(t, eerror.WrapMessageCap, wErr.Cap())

		eerror.DefaultWrapCap(12)
		err = eerror.New("orig err").
			Wrap("[first]").
			Wrap("[second]").
			Err()
		wErr, ok = eerror.Cause(err).(*eerror.WrapError)
		assert.True(t, ok)
		assert.Equal(t, 12, wErr.Cap())
	})
}
