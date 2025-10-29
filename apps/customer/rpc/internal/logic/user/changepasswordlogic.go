package userlogic

import (
	"context"

	"go-dango/apps/customer/rpc/internal/svc"
	"go-dango/apps/customer/rpc/pb"
	"go-dango/pkg/auth"
	"go-dango/pkg/database"
	"go-dango/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(in *pb.ChangePasswordRequest) (*pb.NilOut, error) {
	// todo: add your logic here and delete this line
	if in.NewPassword != in.ConfirmPassword {
		return nil, ErrConfirmPasswordMismatch
	}
	if t := GetPasswordStrength(in.NewPassword); t < StrengthStrong {
		return nil, ErrPasswordStrengthFailed
	}
	uc, rErr := auth.GetUserClaims(l.ctx)
	if rErr != nil {
		l.Logger.Errorw("获取上下文用户信息失败", logx.Field(errors.ErrKey, rErr))
		return nil, rErr
	}
	m, err := l.svcCtx.User.FindModel(l.ctx, []string{}, uc.UserId)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	ok, err := l.svcCtx.Hasher.Verify(in.OldPassword, m.Password)
	if err != nil || !ok {
		return nil, ErrPasswordMismatch
	}
	password, err := l.svcCtx.Hasher.Hash(in.NewPassword)
	if err != nil {
		return nil, ErrPasswordHashError.WithCause(err)
	}
	if err := l.svcCtx.User.UpdateModel(
		l.ctx,
		map[string]any{"password": password},
		map[string]any{"id": uc.UserId},
	); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	return &pb.NilOut{}, nil
}
