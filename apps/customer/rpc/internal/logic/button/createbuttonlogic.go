package buttonlogic

import (
	"context"

	"go-dango/apps/customer/rpc/internal/converter"
	"go-dango/apps/customer/rpc/internal/models"
	"go-dango/apps/customer/rpc/internal/svc"
	"go-dango/apps/customer/rpc/pb"
	"go-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateButtonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateButtonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateButtonLogic {
	return &CreateButtonLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateButtonLogic) CreateButton(in *pb.CreateButtonRequest) (*pb.ButtonOut, error) {
	// todo: add your logic here and delete this line
	mm, err := l.svcCtx.Menu.FindModel(l.ctx, nil, in.MenuId)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	m := models.ButtonModel{
		StandardModel: database.StandardModel{
			BaseModel: database.BaseModel{Id: in.Id},
		},
		Name:         in.Name,
		ArrangeOrder: in.ArrangeOrder,
		IsActive:     in.IsActive,
		Descr:        in.Descr,
		MenuId:       in.MenuId,
		Menu:         *mm,
	}
	pms, err := l.svcCtx.Perm.ListModelByIds(l.ctx, in.PermissionIds)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if len(pms) > 0 {
		m.Permissions = pms
	}
	if err := l.svcCtx.Button.CreateModel(l.ctx, &m); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Button.AddGroupPolicy(l.ctx, m); err != nil {
		return nil, ErrAddButtonPolicy.WithCause(err)
	}
	return converter.ButtonModelToOut(m), nil
}
