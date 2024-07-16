package eerror_test

import (
	"errors"
	"github.com/OddEer0/Eer0/eerror"
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserReq struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Option struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func TestError(t *testing.T) {
	t.Run("Should correct New", func(t *testing.T) {
		err := eerror.New("orig err").Err()
		assert.Equal(t, "orig err", err.Error())
		err = eerror.Cause(err)
		assert.Equal(t, "orig err", err.Error())
	})

	t.Run("Should correct Err", func(t *testing.T) {
		t.Run("Should correct Err with Error", func(t *testing.T) {
			err := eerror.New("orig err").
				Code(404).
				Err()
			eerror.Err(err).
				Code(400)

			assert.Equal(t, eerror.ErrorCode(400), eerror.Status(err))
		})
	})

	t.Run("Should correct Stack", func(t *testing.T) {
		t.Run("Get nil with another err", func(t *testing.T) {
			stack := eerror.GetStack(errors.New("another"))
			assert.Nil(t, stack)
		})

		t.Run("Get empty stack equal nil", func(t *testing.T) {
			err := eerror.New("orig err").Err()
			stack := eerror.GetStack(err)
			assert.Nil(t, stack)
		})

		t.Run("Get correct stack with New", func(t *testing.T) {
			err := eerror.New("orig err").
				Stack("first").
				Stack("second").
				Stack("third").Err()

			stack := eerror.GetStack(err)
			assert.Equal(t, eerror.Stack{"third", "second", "first"}, stack)
		})

		t.Run("Get correct stack with Err", func(t *testing.T) {
			err := eerror.New("orig err").
				Stack("first").
				Err()

			err = eerror.Err(err).
				Stack("second").
				Stack("third").Err()

			stack := eerror.GetStack(err)
			assert.Equal(t, eerror.Stack{"third", "second", "first"}, stack)
		})

		t.Run("Stack added cap + 1", func(t *testing.T) {
			err := eerror.New("orig err").
				Stack(0).
				Err()

			errStruct, ok := err.(*eerror.Error)
			assert.True(t, ok)
			assert.Equal(t, 1, errStruct.StackTrace.Len())
			assert.Equal(t, eerror.DefaultStackSize, errStruct.StackTrace.Cap())

			for i := 0; i < eerror.DefaultStackSize+2; i++ {
				err = eerror.Err(err).
					Stack(i + 1).
					Err()
			}

			assert.Equal(t, 1+eerror.DefaultStackSize+2, errStruct.StackTrace.Len())
			assert.Equal(t, eerror.DefaultStackSize*2, errStruct.StackTrace.Cap())

			assert.Equal(t, eerror.Stack{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, eerror.GetStack(err))
		})

		t.Run("Empty stack get nil value", func(t *testing.T) {
			err := eerror.New("orig err").
				Stack("first").
				Err()

			errStruct, ok := err.(*eerror.Error)
			stack := errStruct.StackTrace.Get()
			assert.Equal(t, eerror.Stack{"first"}, stack)

			assert.True(t, ok)
			val := errStruct.StackTrace.Pop()
			assert.Equal(t, "first", val.(string))
			val = errStruct.StackTrace.Pop()
			assert.Nil(t, val)

			stack = errStruct.StackTrace.Get()
			assert.Nil(t, stack)
		})
	})

	t.Run("Should correct builder", func(t *testing.T) {
		baseErr := errors.New("orig err")
		userParam := UserReq{Id: "1", Name: "incorrect"}
		option := Option{Limit: 10, Page: 1, Sort: "asc", Order: "desc"}

		t.Run("Should correct add msg", func(t *testing.T) {
			err := eerror.Err(baseErr).Msg("kek1").Msg("kek2").Err()
			assert.Equal(t, "kek1", err.Error())
		})

		t.Run("Should correct change code", func(t *testing.T) {
			err := eerror.Err(baseErr).Code(100).Err()
			assert.Equal(t, eerror.ErrorCode(100), eerror.Status(err))
		})

		t.Run("Should correct add info field", func(t *testing.T) {
			assert.Nil(t, eerror.Info(baseErr))
			err := eerror.Err(baseErr).Err()
			assert.Nil(t, eerror.Info(err))
			err = eerror.Err(baseErr).
				Field("err_param", map[string]interface{}{
					"user":   userParam,
					"option": option,
				}).
				Field("test", "test").
				Fields(map[string]interface{}{
					"many_field":  "value",
					"many_field2": "value",
				}).
				Err()
			in := eerror.Info(err)
			assert.Equal(t, map[string]interface{}{
				"err_param": map[string]interface{}{
					"user":   userParam,
					"option": option,
				},
				"test":        "test",
				"many_field":  "value",
				"many_field2": "value",
			}, in)

			err = eerror.Err(baseErr).
				Fields(map[string]interface{}{
					"err_param": map[string]interface{}{
						"user":   userParam,
						"option": option,
					},
					"test":        "test",
					"many_field":  "value",
					"many_field2": "value",
				}).
				Err()

			in = eerror.Info(err)
			assert.Equal(t, map[string]interface{}{
				"err_param": map[string]interface{}{
					"user":   userParam,
					"option": option,
				},
				"test":        "test",
				"many_field":  "value",
				"many_field2": "value",
			}, in)
		})

		t.Run("Should correct cause", func(t *testing.T) {
			err := eerror.Cause(errors.New("another"))
			assert.Nil(t, err)
			err = eerror.OrigCause(errors.New("another"))
			assert.Nil(t, err)
			err = eerror.Err(baseErr).
				Err()

			assert.Equal(t, baseErr.Error(), eerror.Cause(err).Error())

			err = eerror.New("kek wait").Err()
			assert.Equal(t, "kek wait", eerror.Cause(err).Error())
		})

		t.Run("Should correct cause wrap", func(t *testing.T) {
			err := eerror.Err(baseErr).
				Wrap("[first]").
				Wrap("[second]").
				Err()

			expectedErrStr := "[second]" + eerror.DefaultWrapSeparator + "[first]" + eerror.DefaultWrapSeparator + baseErr.Error()
			assert.Equal(t, expectedErrStr, eerror.Cause(err).Error())
			origCause := eerror.OrigCause(err)
			assert.Equal(t, baseErr.Error(), origCause.Error())
		})
	})
}
