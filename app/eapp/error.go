package eapp

import (
	"github.com/pkg/errors"
	"strings"
	"sync"
)

var (
	ErrLibConfigToBeNull  = errors.New("lib config is null")
	ErrUserConfigToBeNull = errors.New("user config is null")
)

type multiError struct {
	errs      []error
	separator string
	mu        sync.Mutex
}

func (m *multiError) Add(err ...error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errs = append(m.errs, err...)
}

func (m *multiError) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.errs)
}

func (m *multiError) Error() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	b := strings.Builder{}

	for i, err := range m.errs {
		b.WriteString(err.Error())
		if i < len(m.errs)-1 {
			b.WriteString(m.separator)
		}
	}

	return b.String()
}

func newMultiError(separator string, errs ...error) *multiError {
	res := &multiError{errs: errs, separator: separator, mu: sync.Mutex{}}
	if len(errs) == 0 {
		res.errs = make([]error, 0)
	}
	return res
}
