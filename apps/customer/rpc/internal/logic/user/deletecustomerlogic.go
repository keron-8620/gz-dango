package userlogic

import (
	"context"

	"go-dango/apps/customer/rpc/internal/svc"
	"go-dango/apps/customer/rpc/pb"
	"go-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCustomerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCustomerLogic {
	return &DeleteCustomerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCustomerLogic) DeleteCustomer(in *pb.DeleteUserRequest) (*pb.NilOut, error) {
	// todo: add your logic here and delete this line
	if err := l.svcCtx.User.DeleteModel(l.ctx, in.Pk); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	return &pb.NilOut{}, nil
}
