package userlogic

import (
	"net/http"

	"go-dango/pkg/errors"
)

var (
	RsonConfirmPasswordMismatch = errors.ReasonEnum{
        Reason: "confirm_password_mismatch",
        Msg:    "两次输入的密码不一致",
    }
	RsonPasswordStrengthFailed = errors.ReasonEnum{
		Reason: "password_strength_failed",
		Msg:    "密码强度不够",
	}
	RsonPasswordHashError = errors.ReasonEnum{
		Reason: "password_hash_error",
		Msg:    "密码加密错误",
	}
	RsonPasswordMismatch = errors.ReasonEnum{
		Reason: "password_mismatch",
		Msg:    "密码错误",
	}
	RsonInvalidCredentials = errors.ReasonEnum{
		Reason: "invalid_credentials",
		Msg:    "用户名或密码错误",
	}
	
	RsonUserInActive = errors.ReasonEnum{
		Reason: "user_inactive",
		Msg:    "用户未激活",
	}
)

var (
	ErrConfirmPasswordMismatch = errors.New(
        http.StatusBadRequest,
        RsonConfirmPasswordMismatch.Reason,
        RsonConfirmPasswordMismatch.Msg,
        nil,
    )
	ErrPasswordStrengthFailed = errors.New(
		http.StatusBadRequest,
		RsonPasswordStrengthFailed.Reason,
		RsonPasswordStrengthFailed.Msg,
		nil,
	)
	ErrPasswordHashError = errors.New(
		http.StatusInternalServerError,
		RsonPasswordHashError.Reason,
		RsonPasswordHashError.Msg,
		nil,
	)
	ErrInvalidCredentials = errors.New(
		http.StatusUnauthorized,
		RsonInvalidCredentials.Reason,
		RsonInvalidCredentials.Msg,
		nil,
	)
	ErrPasswordMismatch = errors.New(
		http.StatusUnauthorized,
		RsonPasswordMismatch.Reason,
		RsonPasswordMismatch.Msg,
		nil,
	)
	ErrUserInActive = errors.New(
		http.StatusUnauthorized,
		RsonUserInActive.Reason,
		RsonUserInActive.Msg,
		nil,
	)
)
