package menulogic

import (
	"net/http"

	"gz-dango/pkg/errors"
)

var (
	ErrAddMenuPolicy = errors.New(
		http.StatusInternalServerError,
		"add_menu_policy_failed",
		"添加菜单策略失败",
		nil,
	)
	ErrRemoveMenuPolicy = errors.New(
		http.StatusInternalServerError,
		"remove_menu_policy_failed",
		"删除菜单策略失败",
		nil,
	)
)
