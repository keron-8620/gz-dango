package svc

import (
	"context"
	"time"

	"go-dango/apps/customer/rpc/internal/models"
	"go-dango/pkg/database"
	"go-dango/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UserService struct {
	gormDB *gorm.DB
}

func NewUserService(
	gormDB *gorm.DB,
) *UserService {
	return &UserService{
		gormDB: gormDB,
	}
}

func (s *UserService) CreateModel(ctx context.Context, m *models.UserModel) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := database.DBCreate(ctx, s.gormDB, &models.UserModel{}, m); err != nil {
		logx.WithContext(ctx).Errorw(
			"新增用户模型失败",
			logx.Field("username", m.Username),
			logx.Field("is_active", m.IsActive),
			logx.Field("is_staff", m.IsStaff),
			logx.Field("role_id", m.RoleId),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}

func (s *UserService) UpdateModel(ctx context.Context, data map[string]any, conds ...any) error {
	if err := database.DBUpdate(ctx, s.gormDB, &models.UserModel{}, data, nil, conds...); err != nil {
		fields := mapToLogFields(data)
		fields = append(fields, logx.Field(errors.ErrKey, err))
		logx.WithContext(ctx).Errorw("更新用户模型失败", fields...)
		return err
	}
	return nil
}

func (s *UserService) DeleteModel(ctx context.Context, conds ...any) error {
	if err := database.DBDelete(ctx, s.gormDB, &models.UserModel{}, conds...); err != nil {
		logx.WithContext(ctx).Errorw(
			"删除用户模型失败",
			logx.Field(database.CondsKey, conds),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}

func (s *UserService) FindModel(
	ctx context.Context,
	preloads []string,
	conds ...any,
) (*models.UserModel, error) {
	var m models.UserModel
	if err := database.DBFind(ctx, s.gormDB, preloads, &m, conds...); err != nil {
		logx.WithContext(ctx).Errorw(
			"查询用户模型失败",
			logx.Field(database.CondsKey, conds),
			logx.Field(errors.ErrKey, err),
		)
		return nil, err
	}
	return &m, nil
}

func (s *UserService) ListModel(
	ctx context.Context,
	qp database.QueryParams,
) (int64, []models.UserModel, error) {
	var ms []models.UserModel
	count, err := database.DBList(ctx, s.gormDB, &models.UserModel{}, &ms, qp)
	if err != nil {
		fields := qpToLogFields(qp)
		fields = append(fields, logx.Field(errors.ErrKey, err))
		logx.WithContext(ctx).Errorw("查询用户列表失败", fields...)
		return 0, nil, err
	}
	return count, ms, err
}
