package menulogic

import (
	"context"
	"time"

	"go-dango/apps/customer/rpc/internal/converter"
	"go-dango/apps/customer/rpc/internal/svc"
	"go-dango/apps/customer/rpc/pb"
	"go-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMenuLogic) UpdateMenu(in *pb.UpdateMenuRequest) (*pb.MenuOut, error) {
	// todo: add your logic here and delete this line
	data := map[string]any{
		"update_at":     time.Now(),
		"path":          in.Path,
		"component":     in.Component,
		"name":          in.Name,
		"meta":          in.Meta,
		"label":         in.Label,
		"arrange_order": in.ArrangeOrder,
		"is_active":     in.IsActive,
		"descr":         in.Descr,
	}
	pms, err := l.svcCtx.Perm.ListModelByIds(l.ctx, in.PermissionIds)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Menu.UpdateModel(l.ctx, data, map[string]any{"permissions": pms}); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	m, err := l.svcCtx.Menu.FindModel(l.ctx, []string{"Parent", "Permissions"}, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Menu.RemoveGroupPolicy(l.ctx, *m, false); err != nil {
		return nil, ErrRemoveMenuPolicy.WithCause(err)
	}
	if err := l.svcCtx.Menu.AddGroupPolicy(l.ctx, *m); err != nil {
		return nil, ErrAddMenuPolicy.WithCause(err)
	}
	return converter.MenuModelToOut(*m), nil
}
