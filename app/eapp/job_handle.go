package eapp

import (
	"context"
	"github.com/OddEer0/Eer0/econtainer/eset"
	"github.com/pkg/errors"
	"sync"
)

func (a *app) startJob() error {
	return nil
}

func initJobHandle(ctx context.Context, wg *sync.WaitGroup, errChan chan<- error, key string, job Job, cfg any) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			errChan <- errors.Errorf("{key: %s, panic: %v}", key, r)
		}
	}()
	err := job.Init(ctx, cfg)
	if err != nil {
		errChan <- errors.Errorf("{key: %s, err: %v}", key, err)
	}
}

func closeJobHandle(ctx context.Context, key string, job Job, wg *sync.WaitGroup, errCh chan<- error, completed eset.Set[string]) {
	defer wg.Done()
	defer completed.Add(key)
	defer func() {
		if r := recover(); r != nil {
			errCh <- errors.Errorf("{key: %s, panic: %v}", key, r)
		}
	}()
	err := job.Close(ctx)
	if err != nil {
		errCh <- errors.Errorf("{key: %s, err: %v}", key, err)
	}
}

func (a *app) initJobs(ctx context.Context) error {
	errCh := make(chan error)

	go func(ctx context.Context, cfg any, errCh chan<- error, jobs map[string]Job) {
		wg := &sync.WaitGroup{}
		for key, job := range jobs {
			wg.Add(1)
			go initJobHandle(ctx, wg, errCh, key, job, cfg)
		}

		wg.Wait()
		errCh <- nil
	}(ctx, a.configs.Client, errCh, a.jobs)

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "[app.initJobs] ctx.Done")
	case err := <-errCh:
		if err != nil {
			return errors.Wrap(err, "[app.initJobs] job.Init")
		}
		return nil
	}
}

func (a *app) closeJob(ctx context.Context) error {
	multiErr := newMultiError(", ")
	successCh := make(chan struct{})
	completed := eset.New[string]()

	go func() {
		wg := &sync.WaitGroup{}
		errCh := make(chan error)
		for key, job := range a.jobs {
			wg.Add(1)
			go closeJobHandle(ctx, key, job, wg, errCh, completed)
		}

		go func() {
			for err := range errCh {
				multiErr.Add(err)
			}
		}()

		wg.Wait()
		close(errCh)
		successCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return closeJobErrHandle(ctx, completed, a.jobs, multiErr)
	case <-successCh:
		return closeJobErrHandle(ctx, completed, a.jobs, multiErr)
	}
}

func closeJobErrHandle(ctx context.Context, completed eset.Set[string], jobs map[string]Job, multiErr *multiError) error {
	var err error
	if ctx.Err() != nil {
		for key := range jobs {
			if !completed.Has(key) {
				multiErr.Add(errors.Errorf("{key: %s, err: %v}", key, ctx.Err()))
			}
		}
	}

	if multiErr.Len() > 0 {
		err = errors.Errorf("{%v}", multiErr)
	}

	if err != nil {
		return errors.Wrap(err, "[app.closeJob] close job errors")
	}
	return nil
}

func (a *app) runAllJob(ctx context.Context) error {
	errCh := make(chan error)

	go func() {
		wg := &sync.WaitGroup{}
		for key, job := range a.jobs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						errCh <- errors.Errorf("{key: %s, panic: %v}", key, r)
					}
				}()
				err := job.Run()
				if err != nil {
					errCh <- errors.Errorf("{key: %s, err: %v}", key, err)
				}
			}()
		}
		wg.Wait()
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "[app.runAllJob] ctx.Done")
	case err := <-errCh:
		if err != nil {
			return errors.Wrap(err, "[app.runAllJob] job.Run")
		}
		return nil
	}
}
