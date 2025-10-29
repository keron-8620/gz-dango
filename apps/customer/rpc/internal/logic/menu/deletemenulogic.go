package menulogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMenuLogic) DeleteMenu(in *pb.DeleteMenuRequest) (*pb.NilOut, error) {
	// todo: add your logic here and delete this line
	m, err := l.svcCtx.Menu.FindModel(l.ctx, nil, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Menu.DeleteModel(l.ctx, in.Pk); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Menu.RemoveGroupPolicy(l.ctx, *m, true); err != nil {
		return nil, ErrRemoveMenuPolicy.WithCause(err)
	}
	return &pb.NilOut{}, nil
}
