package userlogic

import (
	"context"

	"go-dango/apps/customer/rpc/internal/converter"
	"go-dango/apps/customer/rpc/internal/svc"
	"go-dango/apps/customer/rpc/pb"
	"go-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCustomerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCustomerLogic {
	return &GetCustomerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCustomerLogic) GetCustomer(in *pb.GetUserRequest) (*pb.UserOut, error) {
	// todo: add your logic here and delete this line

	m, err := l.svcCtx.User.FindModel(l.ctx, []string{"Role"}, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	return converter.UserModelToOut(*m), nil
}
