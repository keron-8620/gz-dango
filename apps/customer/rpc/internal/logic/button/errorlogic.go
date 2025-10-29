package buttonlogic

import (
	"net/http"

	"gz-dango/pkg/errors"
)

var (
	ErrAddButtonPolicy = errors.New(
		http.StatusInternalServerError,
		"add_button_policy_failed",
		"添加按钮策略失败",
		nil,
	)
	ErrRemoveButtonPolicy = errors.New(
		http.StatusInternalServerError,
		"remove_button_policy_failed",
		"删除按钮策略失败",
		nil,
	)
)
