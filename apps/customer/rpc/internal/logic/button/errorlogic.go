package buttonlogic

import (
	"net/http"

	"go-dango/pkg/errors"
)

var (
	RsonAddButtonPolicy = errors.ReasonEnum{
		Reason: "add_button_policy_failed",
		Msg:    "添加按钮策略失败",
	}
	RsonRemoveButtonPolicy = errors.ReasonEnum{
		Reason: "remove_button_policy_failed",
		Msg:    "删除按钮策略失败",
	}
)

var (
	ErrAddButtonPolicy = errors.New(
		http.StatusInternalServerError,
		RsonAddButtonPolicy.Reason,
		RsonAddButtonPolicy.Msg,
		nil,
	)
	ErrRemoveButtonPolicy = errors.New(
		http.StatusInternalServerError,
		RsonRemoveButtonPolicy.Reason,
		RsonRemoveButtonPolicy.Msg,
		nil,
	)
)
