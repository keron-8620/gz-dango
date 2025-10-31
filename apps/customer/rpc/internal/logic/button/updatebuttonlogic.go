package buttonlogic

import (
	"context"
	"time"

	"gz-dango/apps/customer/rpc/internal/converter"
	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/auth"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateButtonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateButtonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateButtonLogic {
	return &UpdateButtonLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateButtonLogic) UpdateButton(in *pb.UpdateButtonRequest) (*pb.ButtonOut, error) {
	// todo: add your logic here and delete this line
	data := map[string]any{
		"update_at":     time.Now(),
		"name":          in.Name,
		"arrange_order": in.ArrangeOrder,
		"is_active":     in.IsActive,
		"descr":         in.Descr,
		"menu_id":       in.MenuId,
	}

	pms, err := l.svcCtx.Perm.ListModelByIds(l.ctx, in.PermissionIds)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Button.UpdateModel(l.ctx, data, map[string]any{"permissions": pms}); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	m, err := l.svcCtx.Button.FindModel(l.ctx, []string{"Menu", "Permissions"}, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Button.RemoveGroupPolicy(l.ctx, *m, false); err != nil {
		return nil, ErrRemoveButtonPolicy.WithCause(err)
	}
	if err := l.svcCtx.Button.AddGroupPolicy(l.ctx, *m); err != nil {
		return nil, ErrAddButtonPolicy.WithCause(err)
	}
	if err := l.svcCtx.NotifyPolicyChange(); err != nil {
		return nil, auth.ErrCasbinSyncFailed.WithCause(err)
	}
	return converter.ButtonModelToOut(*m), nil
}
