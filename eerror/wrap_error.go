package eerror

import (
	"strings"
)

func (w *WrapError) Error() string {
	sl := strings.Builder{}
	for i := len(w.messages) - 1; i >= 0; i-- {
		sl.WriteString(w.messages[i])
		if i > 0 {
			sl.WriteString(wrapSeparator)
		}
	}
	return sl.String()
}

func newWrapError(message string) error {
	sl := make([]string, 0, wrapDepth)
	sl = append(sl, message)
	return &WrapError{
		messages: sl,
	}
}

func (w *WrapError) Wrap(str string) {
	w.messages = append(w.messages, str)
}

func (w *WrapError) Cap() int {
	return cap(w.messages)
}
