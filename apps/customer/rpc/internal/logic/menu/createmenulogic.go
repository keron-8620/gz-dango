package menulogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/converter"
	"gz-dango/apps/customer/rpc/internal/models"
	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/auth"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMenuLogic) CreateMenu(in *pb.CreateMenuRequest) (*pb.MenuOut, error) {
	// todo: add your logic here and delete this line
	m := models.MenuModel{
		StandardModel: database.StandardModel{
			BaseModel: database.BaseModel{Id: in.Id},
		},
		Path:      in.Path,
		Component: in.Component,
		Name:      in.Name,
		Label:     in.Label,
		Meta: models.Meta{
			Icon:  in.Meta.Icon,
			Title: in.Meta.Title,
		},
		ArrangeOrder: in.ArrangeOrder,
		IsActive:     in.IsActive,
		Descr:        in.Descr,
	}
	if in.ParentId != 0 {
		parent, err := l.svcCtx.Menu.FindModel(l.ctx, nil, in.ParentId)
		if err != nil {
			return nil, database.NewGormError(err, nil)
		}
		m.ParentId = &in.ParentId
		m.Parent = parent
	}
	pms, err := l.svcCtx.Perm.ListModelByIds(l.ctx, in.PermissionIds)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if len(pms) > 0 {
		m.Permissions = pms
	}
	if err := l.svcCtx.Menu.CreateModel(l.ctx, &m); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Menu.AddGroupPolicy(l.ctx, m); err != nil {
		return nil, ErrAddMenuPolicy.WithCause(err)
	}
	if err := l.svcCtx.NotifyPolicyChange(); err != nil {
		return nil, auth.ErrCasbinSyncFailed.WithCause(err)
	}
	return converter.MenuModelToOut(m), nil
}
