package rolelogic

import (
	"context"
	"time"

	"go-dango/apps/customer/rpc/internal/converter"
	"go-dango/apps/customer/rpc/internal/svc"
	"go-dango/apps/customer/rpc/pb"
	"go-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRoleLogic) UpdateRole(in *pb.UpdateRoleRequest) (*pb.RoleOut, error) {
	// todo: add your logic here and delete this line
	data := map[string]any{
		"update_at": time.Now(),
		"name":      in.Name,
		"descr":     in.Descr,
	}

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
	upmap := map[string]any{"Permissions": pms, "Menus": mms, "Buttons": bms}
	if err := l.svcCtx.Role.UpdateModel(l.ctx, data, upmap); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	m, err := l.svcCtx.Role.FindModel(l.ctx, []string{"Permissions", "Menus", "Buttons"}, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Role.RemoveGroupPolicy(l.ctx, *m); err != nil {
		return nil, ErrRemoveRolePolicy.WithCause(err)
	}
	if err := l.svcCtx.Role.AddGroupPolicy(l.ctx, *m); err != nil {
		return nil, ErrAddRolePolicy.WithCause(err)
	}
	return converter.RoleModelToOut(*m), nil
}
