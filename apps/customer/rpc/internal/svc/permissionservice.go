package svc

import (
	"context"
	"strconv"
	"time"

	"gz-dango/apps/customer/rpc/internal/models"
	"gz-dango/pkg/auth"
	"gz-dango/pkg/database"
	"gz-dango/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type PermissionService struct {
	gormDB *gorm.DB
	cache  *auth.AuthEnforcer
}

func NewPermissionService(
	gormDB *gorm.DB,
	cache *auth.AuthEnforcer,
) *PermissionService {
	return &PermissionService{
		gormDB: gormDB,
		cache:  cache,
	}
}

func (s *PermissionService) CreateModel(ctx context.Context, m *models.PermissionModel) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := database.DBCreate(ctx, s.gormDB, &models.PermissionModel{}, m); err != nil {
		logx.WithContext(ctx).Errorw(
			"新增权限模型失败",
			logx.Field("id", m.Id),
			logx.Field("url", m.Url),
			logx.Field("method", m.Method),
			logx.Field("label", m.Label),
			logx.Field("descr", m.Descr),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}

func (s *PermissionService) UpdateModel(ctx context.Context, data map[string]any, conds ...any) error {
	if err := database.DBUpdate(ctx, s.gormDB, &models.PermissionModel{}, data, nil, conds...); err != nil {
		fields := database.MapToLogFields(data)
		fields = append(fields, logx.Field(errors.ErrKey, err))
		logx.WithContext(ctx).Errorw("更新权限模型失败", fields...)
		return err
	}
	return nil
}

func (s *PermissionService) DeleteModel(ctx context.Context, conds ...any) error {
	if err := database.DBDelete(ctx, s.gormDB, &models.PermissionModel{}, conds...); err != nil {
		logx.WithContext(ctx).Errorw(
			"删除权限模型失败",
			logx.Field(database.CondsKey, conds),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}

func (s *PermissionService) FindModel(
	ctx context.Context,
	preloads []string,
	conds ...any,
) (*models.PermissionModel, error) {
	var m models.PermissionModel
	if err := database.DBFind(ctx, s.gormDB, preloads, &m, conds...); err != nil {
		logx.WithContext(ctx).Errorw(
			"查询权限模型失败",
			logx.Field(database.CondsKey, conds),
			logx.Field(errors.ErrKey, err),
		)
		return nil, err
	}
	return &m, nil
}

func (s *PermissionService) ListModel(
	ctx context.Context,
	qp database.QueryParams,
) (int64, []models.PermissionModel, error) {
	var ms []models.PermissionModel
	count, err := database.DBList(ctx, s.gormDB, &models.PermissionModel{}, &ms, qp)
	if err != nil {
		fields := database.QPToLogFields(qp)
		fields = append(fields, logx.Field(errors.ErrKey, err))
		logx.WithContext(ctx).Errorw("查询权限列表失败", fields...)
		return 0, nil, err
	}
	return count, ms, err
}

func (s *PermissionService) ListModelByIds(
	ctx context.Context,
	ids []uint32,
) ([]models.PermissionModel, error) {
	if len(ids) == 0 {
		return []models.PermissionModel{}, nil
	}
	qp := database.NewPksQueryParams(ids)
	_, ms, err := s.ListModel(ctx, qp)
	return ms, err
}

func (s *PermissionService) LoadPolicies(ctx context.Context) error {
	qp := database.QueryParams{
		Preloads: []string{},
		Query:    nil,
		OrderBy:  nil,
		Limit:    0,
		Offset:   0,
		IsCount:  false,
	}
	_, ms, err := s.ListModel(ctx, qp)
	if err != nil {
		return err
	}
	for _, m := range ms {
		if err := s.AddPolicy(ctx, m); err != nil {
			return err
		}
	}
	return nil
}

func (s *PermissionService) AddPolicy(
	ctx context.Context,
	m models.PermissionModel,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	sub := permissionModelToSub(m)
	if err := s.cache.AddPolicy(sub, m.Url, m.Method); err != nil {
		logx.WithContext(ctx).Errorw(
			"添加权限策略失败",
			logx.Field(auth.SubKey, sub),
			logx.Field(auth.ObjKey, m.Url),
			logx.Field(auth.ActKey, m.Method),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}

func (s *PermissionService) RemovePolicy(
	ctx context.Context,
	m models.PermissionModel,
	removeInherited bool,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	sub := permissionModelToSub(m)
	if err := s.cache.RemovePolicy(sub, m.Url, m.Method); err != nil {
		logx.WithContext(ctx).Errorw(
			"删除权限策略失败",
			logx.Field(auth.SubKey, sub),
			logx.Field(auth.ObjKey, m.Url),
			logx.Field(auth.ActKey, m.Method),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	if removeInherited {
		if err := s.cache.RemoveGroupPolicy(1, sub); err != nil {
			logx.WithContext(ctx).Errorw(
				"删除权限作为父级策略失败(该策略被其他策略继承)",
				logx.Field(auth.ObjKey, sub),
				logx.Field(errors.ErrKey, err),
			)
			return err
		}
	}
	return nil
}

func permissionModelToSub(m models.PermissionModel) string {
	return strconv.FormatUint(uint64(m.Id), 10)
}
