package svc

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/models"
	"gz-dango/pkg/database"
	"gz-dango/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	
	"gorm.io/gorm"
)

type RecordService struct {
	gormDB *gorm.DB
}

func NewRecordService(
	gormDB *gorm.DB,
) *RecordService {
	return &RecordService{
		gormDB: gormDB,
	}
}

func (s *RecordService) CreateModel(ctx context.Context, m *models.LoginRecordModel) error {
	if err := database.DBCreate(ctx, s.gormDB, &models.LoginRecordModel{}, m); err != nil {
		logx.WithContext(ctx).Errorw(
			"新增用户登录记录模型失败",
			logx.Field("username", m.Username),
			logx.Field("ip_address", m.IPAddress),
			logx.Field("user_agent", m.UserAgent),
			logx.Field("status", m.Status),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}

func (s *RecordService) DeleteModel(ctx context.Context, conds ...any) error {
	if err := database.DBDelete(ctx, s.gormDB, &models.LoginRecordModel{}, conds...); err != nil {
		logx.WithContext(ctx).Errorw(
			"删除用户登录记录模型失败",
			logx.Field(database.CondsKey, conds),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}

func (s *RecordService) FindModel(
	ctx context.Context,
	preloads []string,
	conds ...any,
) (*models.LoginRecordModel, error) {
	var m models.LoginRecordModel
	if err := database.DBFind(ctx, s.gormDB, preloads, &m, conds...); err != nil {
		logx.WithContext(ctx).Errorw(
			"查询用户登录记录模型失败",
			logx.Field(database.CondsKey, conds),
			logx.Field(errors.ErrKey, err),
		)
		return nil, err
	}
	return &m, nil
}

func (s *RecordService) ListModel(
	ctx context.Context,
	qp database.QueryParams,
) (int64, []models.LoginRecordModel, error) {
	var ms []models.LoginRecordModel
	count, err := database.DBList(ctx, s.gormDB, &models.LoginRecordModel{}, &ms, qp)
	if err != nil {
		fields := database.QPToLogFields(qp)
		fields = append(fields, logx.Field(errors.ErrKey, err))
		logx.WithContext(ctx).Errorw("查询用户登录记录列表失败", fields...)
		return 0, nil, err
	}
	return count, ms, err
}
