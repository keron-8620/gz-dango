package menulogic

import (
	"net/http"

	"go-dango/pkg/errors"
)

var (
	RsonAddMenuPolicy = errors.ReasonEnum{
		Reason: "add_menu_policy_failed",
		Msg:    "添加菜单策略失败",
	}
	RsonRemoveMenuPolicy = errors.ReasonEnum{
		Reason: "remove_menu_policy_failed",
		Msg:    "删除菜单策略失败",
	}
)

var (
	ErrAddMenuPolicy = errors.New(
		http.StatusInternalServerError,
		RsonAddMenuPolicy.Reason,
		RsonAddMenuPolicy.Msg,
		nil,
	)
	ErrRemoveMenuPolicy = errors.New(
		http.StatusInternalServerError,
		RsonRemoveMenuPolicy.Reason,
		RsonRemoveMenuPolicy.Msg,
		nil,
	)
)
