package eerror

import (
	"fmt"
	"github.com/OddEer0/Eer0/eerror"
	"github.com/pkg/errors"
	"os"
	"testing"
)

type TraceMethod struct {
	Type    string `json:"type"`
	Method  string `json:"method"`
	Package string `json:"package"`
}

var (
	trace1 = &TraceMethod{
		Type:    "useCase",
		Method:  "GetById",
		Package: "userUseCase",
	}
	trace2 = &TraceMethod{
		Type:    "service",
		Method:  "SaveById",
		Package: "userService",
	}
	trace3 = &TraceMethod{
		Type:    "useCase",
		Method:  "UpdateMarketById",
		Package: "userUseCase",
	}
	trace4 = &TraceMethod{
		Type:    "grpcHandler",
		Method:  "UpdateMarketById",
		Package: "grpcHandlerV1",
	}

	wrapped1 = "[useCase] db.QueryRow"
	wrapped2 = "[service] userRepository.GetById"
	wrapped3 = "[useCase] userService.SaveById"
	wrapped4 = "[grpcHandler] userService.UpdateMarketById"
)

func BenchmarkError(b *testing.B) {
	openPath := "kek"
	_, errBase := os.Open(openPath)

	b.Run("eerror pkg stack", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := eerror.Err(errBase).
				Stack(trace1).
				Code(eerror.ErrInternal).
				Err()
			err = eerror.Err(err).Stack(trace2).Err()
			err = eerror.Err(err).Stack(trace3).Err()
			err = eerror.Err(err).Stack(trace4).Err()

			err.Error()
			eerror.GetStack(err)
		}
	})

	b.Run("pkg/error wrap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := errors.Wrap(errBase, wrapped1)
			err = errors.Wrap(errBase, wrapped2)
			err = errors.Wrap(errBase, wrapped3)
			err = errors.Wrap(errBase, wrapped4)

			err.Error()
		}
	})

	b.Run("pkg/error WithMessage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := errors.WithMessage(errBase, wrapped1)
			err = errors.WithMessage(errBase, wrapped2)
			err = errors.WithMessage(errBase, wrapped3)
			err = errors.WithMessage(errBase, wrapped4)

			err.Error()
		}
	})

	b.Run("fmt.Errorf wrap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := fmt.Errorf("err: %v", errBase)
			err = fmt.Errorf("err: %v", errBase)
			err = fmt.Errorf("err: %v", errBase)
			err = fmt.Errorf("err: %v", errBase)

			err.Error()
		}
	})
}
