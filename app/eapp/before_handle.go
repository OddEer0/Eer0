package eapp

import (
	"context"
	"github.com/OddEer0/Eer0/econtainer/eset"
	"github.com/pkg/errors"
	"sync"
)

func (a *app) beforeHandle(ctx context.Context) error {
	multiErr := newMultiError(", ")
	successCh := make(chan struct{})
	completed := eset.New[string]()

	go func() {
		errCh := make(chan error, len(a.beforeHandlers))
		wg := &sync.WaitGroup{}

		for _, h := range a.beforeHandlers {
			wg.Add(1)
			go beforeHandle(ctx, h, a.configs.Client, completed, wg, errCh)
		}

		go func(errorCh chan error) {
			for err := range errorCh {
				multiErr.Add(err)
				wg.Done()
			}
		}(errCh)

		wg.Wait()
		close(errCh)
		successCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return beforeErrHandler(ctx, completed, a.beforeHandlers, multiErr)
	case <-successCh:
		return beforeErrHandler(ctx, completed, a.beforeHandlers, multiErr)
	}
}

func beforeErrHandler(ctx context.Context, completed eset.Set[string], handlers map[string]BeforeHandler, multiErr *multiError) error {
	var err error
	if ctx.Err() != nil {
		for _, h := range handlers {
			if !completed.Has(h.Key) {
				multiErr.Add(errors.Errorf("{key: %s, err: %v}", h.Key, ctx.Err()))
			}
		}
	}

	if multiErr.Len() > 0 {
		err = errors.Errorf("{%v}", multiErr)
	}

	if err != nil {
		return errors.Wrap(err, "[app.beforeHandle] before handle errors")
	}
	return nil
}

func beforeHandle(ctx context.Context, h BeforeHandler, cfg any, completed eset.Set[string], wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	defer completed.Add(h.Key)
	defer func() {
		if r := recover(); r != nil {
			wg.Add(1)
			errCh <- errors.Errorf("{key: %s, panic: %v}", h.Key, r)
		}
	}()

	err := h.Handler(ctx, cfg)
	if err != nil {
		wg.Add(1)
		errCh <- errors.Errorf("{key: %s, err: %v}", h.Key, err)
	}
}
