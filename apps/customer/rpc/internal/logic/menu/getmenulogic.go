package menulogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/converter"
	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuLogic {
	return &GetMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuLogic) GetMenu(in *pb.GetMenuRequest) (*pb.MenuOut, error) {
	// todo: add your logic here and delete this line
	m, err := l.svcCtx.Menu.FindModel(l.ctx, []string{"Parent", "Permissons"}, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	return converter.MenuModelToOut(*m), nil
}
