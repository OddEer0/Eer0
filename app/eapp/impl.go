package eapp

import (
	"context"
	"time"
)

func (a *App) Err() error {
	return a.err
}

type Option struct {
	UserCfgInterceptor UserConfigInterceptor
	LibCfgInterceptor  LibConfigInterceptor
}

func (a *App) WithOptions(opt *Option) *App {
	if a.err != nil {
		return nil
	}
	a.userCfgInterceptor = opt.UserCfgInterceptor
	a.libCfgInterceptor = opt.LibCfgInterceptor
	return a
}

func (a *App) WithJobs(jobs ...JobOption) *App {
	for _, job := range jobs {
		a.jobs[job.Key] = job.Job
	}
	return a
}

func (a *App) BeforeHandle(handlers ...BeforeHandler) *App {
	for _, h := range handlers {
		a.beforeHandlers[h.Key] = h
	}
	return a
}

func (a *App) AfterHandle(handlers ...AfterHandler) *App {
	for _, h := range handlers {
		a.afterHandlers[h.Key] = h
	}
	return a
}

func (a *App) Start() error {
	if a.err != nil {
		return a.err
	}
	ctxTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*a.configs.Lib.StartTimeout)
	defer cancelFunc()
	err := a.beforeHandle(ctxTimeout)
	if err != nil {
		return err
	}

	err = a.initJobs(ctxTimeout)
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	err = a.runAllJob(ctx)

	stopErr := a.Stop()
	if stopErr != nil {
		panic(stopErr)
	}

	if err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() error {
	ctxTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*a.configs.Lib.StopTimeout)
	defer cancelFunc()
	err := a.closeJob(ctxTimeout)
	if err != nil {
		return err
	}

	err = a.afterHandle(ctxTimeout)
	if err != nil {
		return err
	}

	return nil
}
