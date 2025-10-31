package userlogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/database"
	"gz-dango/pkg/errors"

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
	m, err := l.svcCtx.User.FindModel(l.ctx, nil, "username = ?", in.Username)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	ok, err := hasher.Verify(in.Password, m.Password)
	if err != nil {
		logx.WithContext(l.ctx).Errorw("")
		return nil, errors.FromError(err)
	}
	if !ok {
		return nil, ErrInvalidCredentials
	}

	return &pb.LoginOut{}, nil
}
