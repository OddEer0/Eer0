package eapp

import (
	"context"
	"time"
)

func (a *app) Err() error {
	return a.err
}

func (a *app) WithJobs(jobs ...JobOption) App {
	for _, job := range jobs {
		a.jobs[job.Key] = job.Job
	}
	return a
}

func (a *app) BeforeHandle(handlers ...BeforeHandler) App {
	for _, h := range handlers {
		a.beforeHandlers[h.Key] = h
	}
	return a
}

func (a *app) AfterHandle(handlers ...AfterHandler) App {
	for _, h := range handlers {
		a.afterHandlers[h.Key] = h
	}
	return a
}

func (a *app) Start() error {
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

func (a *app) Stop() error {
	ctxTimeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*a.configs.Lib.StopTimeout)
	defer cancelFunc()
	err := a.closeJob(ctxTimeout)
	if err != nil {
		return err
	}

	// TODO - обработка всех after handlers

	return nil
}
