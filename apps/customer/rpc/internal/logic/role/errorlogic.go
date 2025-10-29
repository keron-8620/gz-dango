package rolelogic

import (
	"net/http"

	"go-dango/pkg/errors"
)

var (
	RsonAddRolePolicy = errors.ReasonEnum{
		Reason: "add_role_policy_failed",
		Msg:    "添加角色策略失败",
	}
	RsonRemoveRolePolicy = errors.ReasonEnum{
		Reason: "remove_role_policy_failed",
		Msg:    "删除角色策略失败",
	}
)

var (
	ErrAddRolePolicy = errors.New(
		http.StatusInternalServerError,
		RsonAddRolePolicy.Reason,
		RsonAddRolePolicy.Msg,
		nil,
	)
	ErrRemoveRolePolicy = errors.New(
		http.StatusInternalServerError,
		RsonRemoveRolePolicy.Reason,
		RsonRemoveRolePolicy.Msg,
		nil,
	)
)
