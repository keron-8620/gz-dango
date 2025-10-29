package permissionlogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/converter"
	"gz-dango/apps/customer/rpc/internal/models"
	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePermissionLogic {
	return &CreatePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePermissionLogic) CreatePermission(in *pb.CreatePermissionRequest) (*pb.PermissionOutBase, error) {
	// todo: add your logic here and delete this line
	m := models.PermissionModel{
		StandardModel: database.StandardModel{
			BaseModel: database.BaseModel{Id: in.Id},
		},
		Url:    in.Url,
		Method: in.Method,
		Label:  in.Label,
		Descr:  in.Descr,
	}
	if err := l.svcCtx.Perm.CreateModel(l.ctx, &m); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	if err := l.svcCtx.Perm.AddPolicy(l.ctx, m); err != nil {
		return nil, ErrAddPermissionPolicy.WithCause(err)
	}
	return converter.PermModelToOutBase(m), nil
}
