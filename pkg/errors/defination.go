package errors

// import (
// 	"context"
// 	"errors"
// 	"net/http"
// )

// const ErrKey = "error"

// func FromError(err error) *Error {
// 	if err == nil {
// 		return nil
// 	}
// 	if se := new(Error); errors.As(err, &se) {
// 		return se
// 	}
// 	if errors.Is(err, context.Canceled) {
// 		return ErrCtxCancel.WithCause(err)
// 	}
// 	if errors.Is(err, context.DeadlineExceeded) {
// 		return ErrCtxDeadlineExceeded.WithCause(err)
// 	}
// 	return ErrUnKnow.WithCause(err)
// }

// var (
// 	ErrUnKnow = New(
// 		http.StatusInternalServerError,
// 		RsonUnknown.Reason,
// 		RsonUnknown.Msg,
// 		nil,
// 	)
// 	ErrCtxCancel = New(
// 		http.StatusBadRequest,
// 		RsonCtxCancel.Reason,
// 		RsonCtxCancel.Msg,
// 		nil,
// 	)
// 	ErrCtxDeadlineExceeded = New(
// 		http.StatusBadRequest,
// 		RsonCtxDeadlineExceeded.Reason,
// 		RsonCtxDeadlineExceeded.Msg,
// 		nil,
// 	)
// )
