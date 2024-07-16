package eerror

import (
	"errors"
	"github.com/OddEer0/Eer0/eerror"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultErrorMethod(t *testing.T) {
	baseErr := errors.New("baseErr")

	err := eerror.Err(baseErr).Code(eerror.ErrBadRequest).Err()
	assert.Equal(t, eerror.MsgBadRequest, err.Error())
	err = eerror.Err(baseErr).Code(eerror.ErrUnauthorized).Err()
	assert.Equal(t, eerror.MsgUnauthorized, err.Error())
	err = eerror.Err(baseErr).Code(eerror.ErrForbidden).Err()
	assert.Equal(t, eerror.MsgForbidden, err.Error())
	err = eerror.Err(baseErr).Code(eerror.ErrNotFound).Err()
	assert.Equal(t, eerror.MsgNotFound, err.Error())
	err = eerror.Err(baseErr).Code(eerror.ErrDeadlineExceeded).Err()
	assert.Equal(t, eerror.MsgDeadlineExceeded, err.Error())
	err = eerror.Err(baseErr).Code(eerror.ErrConflict).Err()
	assert.Equal(t, eerror.MsgConflict, err.Error())
	err = eerror.Err(baseErr).Code(eerror.ErrUnprocessable).Err()
	assert.Equal(t, eerror.MsgUnprocessable, err.Error())
	err = eerror.Err(baseErr).Code(eerror.ErrInternal).Err()
	assert.Equal(t, eerror.MsgInternal, err.Error())
	err = eerror.Err(baseErr).Code(eerror.ErrUnimplemented).Err()
	assert.Equal(t, eerror.MsgUnimplemented, err.Error())
	err = eerror.Err(baseErr).Code(eerror.ErrBadGateway).Err()
	assert.Equal(t, eerror.MsgBadGateway, err.Error())
	err = eerror.Err(baseErr).Code(eerror.ErrServiceUnavailable).Err()
	assert.Equal(t, eerror.MsgServiceUnavailable, err.Error())
	err = eerror.Err(baseErr).Code(eerror.ErrGatewayTimeout).Err()
	assert.Equal(t, eerror.MsgGatewayTimeout, err.Error())
}
