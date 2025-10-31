package permissionlogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/auth"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePermissionLogic {
	return &DeletePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePermissionLogic) DeletePermission(in *pb.DeletePermissionRequest) (*pb.NilOut, error) {
	// todo: add your logic here and delete this line
	m, err := l.svcCtx.Perm.FindModel(l.ctx, nil, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Perm.DeleteModel(l.ctx, in.Pk); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Perm.RemovePolicy(l.ctx, *m, true); err != nil {
		return nil, ErrRemovePermissionPolicy.WithCause(err)
	}
	if err := l.svcCtx.NotifyPolicyChange(); err != nil {
		return nil, auth.ErrCasbinSyncFailed.WithCause(err)
	}
	return &pb.NilOut{}, nil
}
