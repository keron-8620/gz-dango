package userlogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/converter"
	"gz-dango/apps/customer/rpc/internal/models"
	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *pb.CreateUserRequest) (*pb.UserOut, error) {
	// todo: add your logic here and delete this line
	if t := GetPasswordStrength(in.Password); t < StrengthStrong {
		return nil, ErrPasswordStrengthFailed
	}
	password, err := hasher.Hash(in.Password)
	if err != nil {
		return nil, ErrPasswordHashError.WithCause(err)
	}
	rm, err := l.svcCtx.Role.FindModel(l.ctx, nil, in.RoleId)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	m := models.UserModel{
		Username: in.Username,
		Password: password,
		IsActive: in.IsActive,
		IsStaff:  in.IsStaff,
		RoleId:   in.RoleId,
		Role:     *rm,
	}
	if err := l.svcCtx.User.CreateModel(l.ctx, &m); err != nil {
		return nil, database.NewGormError(err, nil)
	}
	return converter.UserModelToOut(m), nil
}
