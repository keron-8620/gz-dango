package buttonlogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/auth"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteButtonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteButtonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteButtonLogic {
	return &DeleteButtonLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteButtonLogic) DeleteButton(in *pb.DeleteButtonRequest) (*pb.NilOut, error) {
	// todo: add your logic here and delete this line
	m, err := l.svcCtx.Button.FindModel(l.ctx, nil, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Button.DeleteModel(l.ctx, in.Pk); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Button.RemoveGroupPolicy(l.ctx, *m, true); err != nil {
		return nil, ErrRemoveButtonPolicy.WithCause(err)
	}
	if err := l.svcCtx.NotifyPolicyChange(); err != nil {
		return nil, auth.ErrCasbinSyncFailed.WithCause(err)
	}
	return &pb.NilOut{}, nil
}
