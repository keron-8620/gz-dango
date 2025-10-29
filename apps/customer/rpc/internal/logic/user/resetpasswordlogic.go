package userlogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(in *pb.ResetPasswordRequest) (*pb.NilOut, error) {
	// todo: add your logic here and delete this line
	if t := GetPasswordStrength(in.Password); t < StrengthStrong {
		return nil, ErrPasswordStrengthFailed
	}
	password, err := l.svcCtx.Hasher.Hash(in.Password)
	if err != nil {
		return nil, ErrPasswordHashError.WithCause(err)
	}
	if err := l.svcCtx.User.UpdateModel(l.ctx, map[string]any{"password": password}, map[string]any{"id": in.Pk}); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	return &pb.NilOut{}, nil
}
