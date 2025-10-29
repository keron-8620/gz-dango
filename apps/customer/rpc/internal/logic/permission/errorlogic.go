package permissionlogic

import (
	"net/http"

	"go-dango/pkg/errors"
)

var (
	RsonAddPermissionPolicy = errors.ReasonEnum{
        Reason: "add_permission_policy_failed",
        Msg:    "添加权限策略失败",
    }
	RsonRemovePermissionPolicy = errors.ReasonEnum{
        Reason: "remove_permission_policy_failed",
        Msg:    "删除权限策略失败",
    }
)

var (
	ErrAddPermissionPolicy = errors.New(
        http.StatusInternalServerError,
        RsonAddPermissionPolicy.Reason,
        RsonAddPermissionPolicy.Msg,
        nil,
    )
	ErrRemovePermissionPolicy = errors.New(
        http.StatusInternalServerError,
        RsonRemovePermissionPolicy.Reason,
        RsonRemovePermissionPolicy.Msg,
        nil,
    )
)
