package eapp

import "context"

type (
	Job interface {
		Init(ctx context.Context, config any) error
		Run() error
		Close(ctx context.Context) error
	}

	JobOption struct {
		Key string
		Job Job
	}

	BeforeHandler struct {
		Key     string
		Handler func(ctx context.Context, config any) error
	}

	AfterHandler struct {
		Key     string
		Handler func(ctx context.Context, config any) error
	}

	App interface {
		LibConfig(conf *LibConfig) App
		LibConfigByYaml(filePath string) App
		WithConfig(conf any) App
		BeforeHandle(handlers ...BeforeHandler) App
		AfterHandle(handlers ...AfterHandler) App
		WithJobs(jobs ...JobOption) App

		Start() error
		Stop() error

		Err() error
	}
)
