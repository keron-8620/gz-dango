package rolelogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteRoleLogic) DeleteRole(in *pb.DeleteRoleRequest) (*pb.NilOut, error) {
	// todo: add your logic here and delete this line
	m, err := l.svcCtx.Role.FindModel(l.ctx, nil, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Role.DeleteModel(l.ctx, in.Pk); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Role.RemoveGroupPolicy(l.ctx, *m); err != nil {
		return nil, ErrRemoveRolePolicy.WithCause(err)
	}
	return &pb.NilOut{}, nil
}
