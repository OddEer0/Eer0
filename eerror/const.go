package eerror

const (
	DefaultErrorCode     ErrorCode = ErrInternal
	DefaultStackSize               = 8
	DefaultWrapSeparator           = ": "
	WrapMessageCap                 = 8
	DefaultInfoCap                 = 16

	NotCode               ErrorCode = 0
	ErrBadRequest         ErrorCode = 400
	ErrUnauthorized       ErrorCode = 401
	ErrForbidden          ErrorCode = 403
	ErrNotFound           ErrorCode = 404
	ErrDeadlineExceeded   ErrorCode = 408
	ErrConflict           ErrorCode = 409
	ErrUnprocessable      ErrorCode = 422
	ErrInternal           ErrorCode = 500
	ErrUnimplemented      ErrorCode = 501
	ErrBadGateway         ErrorCode = 502
	ErrServiceUnavailable ErrorCode = 503
	ErrGatewayTimeout     ErrorCode = 504

	MsgBadRequest         = "Bad request"
	MsgUnauthorized       = "Unauthorized"
	MsgForbidden          = "Forbidden"
	MsgNotFound           = "Not found"
	MsgDeadlineExceeded   = "Deadline exceeded"
	MsgConflict           = "Conflict"
	MsgUnprocessable      = "Unprocessable"
	MsgInternal           = "Internal"
	MsgUnimplemented      = "Unimplemented"
	MsgBadGateway         = "Bad Gateway"
	MsgServiceUnavailable = "Service Unavailable"
	MsgGatewayTimeout     = "Gateway Timeout"
)

var (
	DefaultErrorMethodFunc ErrorMethod = func(err *Error) string {
		if err.Message == "" {
			switch err.StatusCode {
			case ErrBadRequest:
				return MsgBadRequest
			case ErrUnauthorized:
				return MsgUnauthorized
			case ErrForbidden:
				return MsgForbidden
			case ErrNotFound:
				return MsgNotFound
			case ErrDeadlineExceeded:
				return MsgDeadlineExceeded
			case ErrConflict:
				return MsgConflict
			case ErrUnprocessable:
				return MsgUnprocessable
			case ErrInternal:
				return MsgInternal
			case ErrUnimplemented:
				return MsgUnimplemented
			case ErrBadGateway:
				return MsgBadGateway
			case ErrServiceUnavailable:
				return MsgServiceUnavailable
			case ErrGatewayTimeout:
				return MsgGatewayTimeout
			}
		}
		return err.Message
	}
)
