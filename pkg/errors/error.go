package errors

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
)

const ErrKey = "error"

type Error struct {
	Code   int            `json:"code"`
	Reason string         `json:"reason"`
	Msg    string         `json:"msg"`
	Data   map[string]any `json:"data"`
	cause  error
}

func New(code int, reason, message string, data map[string]any) *Error {
	return &Error{
		Code:   code,
		Reason: reason,
		Msg:    message,
		Data:   data,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code=%d reason = %s msg = %s data = %v cause = %v", e.Code, e.Reason, e.Msg, e.Data, e.cause)
}

func (e *Error) Unwrap() error { return e.cause }

func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Code == e.Code && se.Reason == e.Reason
	}
	return false
}

func (e *Error) WithCause(cause error) *Error {
	err := Clone(e)
	err.cause = cause
	return err
}

func (e *Error) WithData(md map[string]any) *Error {
	err := Clone(e)
	err.Data = md
	return err
}

func (e *Error) Reply() map[string]any {
	data := e.Data
	if e.cause != nil {
		data = make(map[string]any, len(e.Data)+1)
		maps.Copy(data, e.Data)
		data["cause"] = e.cause.Error()
	}
	if data == nil {
		data = map[string]any{}
	}
	return map[string]any{
		"code":   e.Code,
		"reason": e.Reason,
		"msg":    e.Msg,
		"data":   data,
	}
}

func Clone(err *Error) *Error {
	if err == nil {
		return nil
	}
	metadata := make(map[string]any, len(err.Data))
	maps.Copy(metadata, err.Data)
	return &Error{
		Code:   err.Code,
		Reason: err.Reason,
		Msg:    err.Msg,
		Data:   metadata,
		cause:  err.cause,
	}
}

const (
	unknownCode   = http.StatusInternalServerError
	unknownReason = "unknown"
	unknownMsg    = "未知错误"

	ctxCancelCode   = http.StatusBadRequest
	ctxCancelReason = "ctx_cancel"
	ctxCancelMsg    = "请求取消"

	ctxDeadlineCode   = http.StatusRequestTimeout
	ctxDeadlineReason = "ctx_deadline"
	ctxDeadlineMsg    = "请求超时"
)

func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	if errors.Is(err, context.Canceled) {
		return &Error{
			Code:   ctxCancelCode,
			Reason: ctxCancelReason,
			Msg:    ctxCancelMsg,
			Data:   nil,
			cause:  err,
		}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return &Error{
			Code:   ctxDeadlineCode,
			Reason: ctxDeadlineReason,
			Msg:    ctxDeadlineMsg,
			Data:   nil,
			cause:  err,
		}
	}
	return &Error{
		Code:   unknownCode,
		Reason: unknownReason,
		Msg:    unknownMsg,
		Data:   nil,
		cause:  err,
	}
}
