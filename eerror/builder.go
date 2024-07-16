package eerror

func (e *Error) Err() error {
	return e
}

func (e *Error) Stack(item interface{}) ErrorBuilder {
	if e.StackTrace == nil {
		e.StackTrace = newReverseSlice(stackDepth)
	}
	e.StackTrace.Push(item)
	return e
}

func (e *Error) Code(code ErrorCode) ErrorBuilder {
	e.StatusCode = code
	return e
}

func (e *Error) Msg(message string) ErrorBuilder {
	if e.Message == "" {
		e.Message = message
	}
	return e
}

func (e *Error) Field(key string, value interface{}) ErrorBuilder {
	if e.Info == nil {
		e.Info = make(map[string]interface{}, infoCap)
	}
	if _, has := e.Info[key]; !has {
		e.Info[key] = value
	}
	return e
}

func (e *Error) Fields(fields map[string]interface{}) ErrorBuilder {
	if e.Info == nil {
		e.Info = make(map[string]interface{}, infoCap)
	}
	for key, value := range fields {
		e.Info[key] = value
	}
	return e
}

func (e *Error) Wrap(str string) ErrorBuilder {
	bErr, ok := e.Cause.(*WrapError)
	if ok {
		bErr.Wrap(str)
		return e
	}
	e.Cause = newWrapError(e.Cause.Error())
	bErr = e.Cause.(*WrapError)
	bErr.Wrap(str)
	return e
}
