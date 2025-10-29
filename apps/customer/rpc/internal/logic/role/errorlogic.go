package rolelogic

import (
	"net/http"

	"gz-dango/pkg/errors"
)

var (
	ErrAddRolePolicy = errors.New(
		http.StatusInternalServerError,
		"add_role_policy_failed",
		"添加角色策略失败",
		nil,
	)
	ErrRemoveRolePolicy = errors.New(
		http.StatusInternalServerError,
		"remove_role_policy_failed",
		"删除角色策略失败",
		nil,
	)
)
