package userlogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginRequest) (*pb.LoginOut, error) {
	// todo: add your logic here and delete this line

	return &pb.LoginOut{}, nil
}
