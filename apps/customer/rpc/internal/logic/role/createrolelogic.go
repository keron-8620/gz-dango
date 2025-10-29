package rolelogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/converter"
	"gz-dango/apps/customer/rpc/internal/models"
	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRoleLogic) CreateRole(in *pb.CreateRoleRequest) (*pb.RoleOut, error) {
	// todo: add your logic here and delete this line
	pms, err := l.svcCtx.Perm.ListModelByIds(l.ctx, in.PermissionIds)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	mms, err := l.svcCtx.Menu.ListModelByIds(l.ctx, in.MenuIds)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	bms, err := l.svcCtx.Button.ListModelByIds(l.ctx, in.ButtonIds)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	m := models.RoleModel{
		Name:        in.Name,
		Descr:       in.Descr,
		Permissions: pms,
		Menus:       mms,
		Buttons:     bms,
	}
	if err := l.svcCtx.Role.CreateModel(l.ctx, &m); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Role.AddGroupPolicy(l.ctx, m); err != nil {
		return nil, ErrAddRolePolicy.WithCause(err)
	}
	return converter.RoleModelToOut(m), nil
}
