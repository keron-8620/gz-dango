package rolelogic

import (
	"context"

	"go-dango/apps/customer/rpc/internal/converter"
	"go-dango/apps/customer/rpc/internal/svc"
	"go-dango/apps/customer/rpc/pb"
	"go-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRoleLogic) GetRole(in *pb.GetRoleRequest) (*pb.RoleOut, error) {
	// todo: add your logic here and delete this line

	m, err := l.svcCtx.Role.FindModel(l.ctx, []string{"Permissons", "Menus", "Buttons"}, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	return converter.RoleModelToOut(*m), nil
}
