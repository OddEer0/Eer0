package eapp

import (
	"github.com/OddEer0/Eer0/etools/struct/estruct"
	"github.com/pkg/errors"
)

type (
	InitOptions struct {
		Name        string
		Version     string
		Description string
	}
)

func Init(opt *InitOptions) *App {
	err := estruct.NotZero(*opt)
	result := &App{}
	if err != nil {
		result.err = errors.Wrap(err, "[eapp] Init")
		return result
	}

	// init struct
	result.info = Info{
		name:        opt.Name,
		version:     opt.Version,
		description: opt.Description,
	}
	result.configs = &Configs{}
	result.beforeHandlers = make(map[string]BeforeHandler, 8)
	result.afterHandlers = make(map[string]AfterHandler, 8)
	result.jobs = make(map[string]Job, 8)

	return result
}
