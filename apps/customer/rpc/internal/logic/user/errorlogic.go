package userlogic

import (
	"net/http"

	"gz-dango/pkg/errors"
)

var (
	ErrConfirmPasswordMismatch = errors.New(
		http.StatusBadRequest,
		"confirm_password_mismatch",
		"两次输入的密码不一致",
		nil,
	)
	ErrPasswordStrengthFailed = errors.New(
		http.StatusBadRequest,
		"password_strength_failed",
		"密码强度不够",
		nil,
	)
	ErrPasswordHashError = errors.New(
		http.StatusInternalServerError,
		"password_hash_error",
		"密码加密错误",
		nil,
	)
	ErrInvalidCredentials = errors.New(
		http.StatusUnauthorized,
		"invalid_credentials",
		"用户名或密码错误",
		nil,
	)
	ErrPasswordMismatch = errors.New(
		http.StatusUnauthorized,
		"password_mismatch",
		"密码错误",
		nil,
	)
	ErrUserInActive = errors.New(
		http.StatusUnauthorized,
		"user_inactive",
		"用户未激活",
		nil,
	)
)
