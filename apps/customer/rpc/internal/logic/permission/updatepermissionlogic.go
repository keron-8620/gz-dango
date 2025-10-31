package permissionlogic

import (
	"context"
	"time"

	"gz-dango/apps/customer/rpc/internal/converter"
	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/auth"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePermissionLogic {
	return &UpdatePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePermissionLogic) UpdatePermission(in *pb.UpdatePermissionRequest) (*pb.PermissionOutBase, error) {
	// todo: add your logic here and delete this line

	data := map[string]any{
		"update_at": time.Now(),
		"url":       in.Url,
		"method":    in.Method,
		"label":     in.Label,
		"descr":     in.Descr,
	}
	if err := l.svcCtx.Perm.UpdateModel(l.ctx, data); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	m, err := l.svcCtx.Perm.FindModel(l.ctx, nil, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Perm.RemovePolicy(l.ctx, *m, false); err != nil {
		return nil, ErrRemovePermissionPolicy.WithCause(err)
	}
	if err := l.svcCtx.Perm.AddPolicy(l.ctx, *m); err != nil {
		return nil, ErrAddPermissionPolicy.WithCause(err)
	}
	if err := l.svcCtx.NotifyPolicyChange(); err != nil {
		return nil, auth.ErrCasbinSyncFailed.WithCause(err)
	}
	return converter.PermModelToOutBase(*m), nil
}
