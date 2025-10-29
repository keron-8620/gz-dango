package userlogic

import (
	"context"
	"time"

	"go-dango/apps/customer/rpc/internal/converter"
	"go-dango/apps/customer/rpc/internal/svc"
	"go-dango/apps/customer/rpc/pb"
	"go-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCustomerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCustomerLogic {
	return &UpdateCustomerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCustomerLogic) UpdateCustomer(in *pb.UpdateUserRequest) (*pb.UserOut, error) {
	// todo: add your logic here and delete this line
	data := map[string]any{
		"update_at": time.Now(),
		"username":  in.Username,
		"is_active": in.IsActive,
		"is_staff":  in.IsStaff,
		"role_id":   in.RoleId,
	}
	if err := l.svcCtx.User.UpdateModel(l.ctx, data, in.Pk); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	m, err := l.svcCtx.User.FindModel(l.ctx, []string{"Role"}, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	return converter.UserModelToOut(*m), nil
}
