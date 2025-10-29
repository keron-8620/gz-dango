package permissionlogic

import (
	"context"

	"go-dango/apps/customer/rpc/internal/converter"
	"go-dango/apps/customer/rpc/internal/svc"
	"go-dango/apps/customer/rpc/pb"
	"go-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPermissionLogic {
	return &GetPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPermissionLogic) GetPermission(in *pb.GetPermissionRequest) (*pb.PermissionOutBase, error) {
	// todo: add your logic here and delete this line
	m, err := l.svcCtx.Perm.FindModel(l.ctx, nil, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	return converter.PermModelToOutBase(*m), nil
}
