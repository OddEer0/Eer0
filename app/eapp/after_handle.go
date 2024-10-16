package eapp

import (
	"context"
	"fmt"
	"github.com/OddEer0/Eer0/econtainer/eset"
	"github.com/pkg/errors"
	"sync"
)

func afterHandleFunc(ctx context.Context, wg *sync.WaitGroup, errCh chan<- error, handler AfterHandler, cfg any, complete eset.Set[string]) {
	defer wg.Done()
	defer complete.Add(handler.Key)
	defer func() {
		if err := recover(); err != nil {
			wg.Add(1)
			errCh <- errors.Errorf("{key: %s, panic: %v}", handler.Key, err)
		}
	}()

	err := handler.Handler(ctx, cfg)
	if err != nil {
		wg.Add(1)
		errCh <- errors.Errorf("{key: %s, err: %s}", handler.Key, err)
	}
}

func (a *app) afterHandle(ctx context.Context) error {
	complete := eset.New[string]()
	successCh := make(chan struct{})
	multiErr := newMultiError(", ")

	go func() {
		wg := &sync.WaitGroup{}
		errCh := make(chan error, len(a.afterHandlers))

		for _, h := range a.afterHandlers {
			wg.Add(1)
			go afterHandleFunc(ctx, wg, errCh, h, a.configs.Client, complete)
		}

		go func(errCh <-chan error) {
			for err := range errCh {
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
		return afterErrHandler(ctx, complete, a.afterHandlers, multiErr)
	case <-successCh:
		return afterErrHandler(ctx, complete, a.afterHandlers, multiErr)
	}
}

func afterErrHandler(ctx context.Context, completed eset.Set[string], handlers map[string]AfterHandler, multiErr *multiError) error {
	var err error
	if ctx.Err() != nil {
		for _, h := range handlers {
			if !completed.Has(h.Key) {
				multiErr.Add(errors.Errorf("{key: %s, err: %v}", h.Key, ctx.Err()))
			}
		}
	}

	if multiErr.Len() > 0 {
		fmt.Println(multiErr)
		err = errors.Errorf("{%v}", multiErr)
	}

	if err != nil {
		return errors.Wrap(err, "[app.beforeHandle] after handle errors")
	}
	return nil
}
