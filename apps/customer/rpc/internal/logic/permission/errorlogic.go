package permissionlogic

import (
	"net/http"

	"gz-dango/pkg/errors"
)

var (
	ErrAddPermissionPolicy = errors.New(
		http.StatusInternalServerError,
		"add_permission_policy_failed",
		"添加权限策略失败",
		nil,
	)
	ErrRemovePermissionPolicy = errors.New(
		http.StatusInternalServerError,
		"remove_permission_policy_failed",
		"删除权限策略失败",
		nil,
	)
)
