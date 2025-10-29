package svc

import (
	"context"
	"fmt"
	"time"

	"gz-dango/apps/customer/rpc/internal/models"
	"gz-dango/pkg/auth"
	"gz-dango/pkg/database"
	"gz-dango/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RoleService struct {
	gormDB *gorm.DB
	cache  *auth.AuthEnforcer
}

func NewRoleService(
	gormDB *gorm.DB,
	cache *auth.AuthEnforcer,
) *RoleService {
	return &RoleService{
		gormDB: gormDB,
		cache:  cache,
	}
}

func (s *RoleService) CreateModel(ctx context.Context, m *models.RoleModel) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := database.DBCreate(ctx, s.gormDB, &models.RoleModel{}, m); err != nil {
		logx.WithContext(ctx).Errorw(
			"新增角色模型失败",
			logx.Field("name", m.Name),
			logx.Field("descr", m.Descr),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}

func (s *RoleService) UpdateModel(ctx context.Context, data map[string]any, upmap map[string]any, conds ...any) error {
	if err := database.DBUpdate(ctx, s.gormDB, &models.RoleModel{}, data, upmap, conds...); err != nil {
		fields := mapToLogFields(data)
		fields = append(fields, logx.Field(errors.ErrKey, err))
		logx.WithContext(ctx).Errorw("更新角色模型失败", fields...)
		return err
	}
	return nil
}

func (s *RoleService) DeleteModel(ctx context.Context, conds ...any) error {
	if err := database.DBDelete(ctx, s.gormDB, &models.RoleModel{}, conds...); err != nil {
		logx.WithContext(ctx).Errorw(
			"删除角色模型失败",
			logx.Field(database.CondsKey, conds),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}

func (s *RoleService) FindModel(
	ctx context.Context,
	preloads []string,
	conds ...any,
) (*models.RoleModel, error) {
	var m models.RoleModel
	if err := database.DBFind(ctx, s.gormDB, preloads, &m, conds...); err != nil {
		logx.WithContext(ctx).Errorw(
			"查询角色模型失败",
			logx.Field(database.CondsKey, conds),
			logx.Field(errors.ErrKey, err),
		)
		return nil, err
	}
	return &m, nil
}

func (s *RoleService) ListModel(
	ctx context.Context,
	qp database.QueryParams,
) (int64, []models.RoleModel, error) {
	var ms []models.RoleModel
	count, err := database.DBList(ctx, s.gormDB, &models.RoleModel{}, &ms, qp)
	if err != nil {
		fields := qpToLogFields(qp)
		fields = append(fields, logx.Field(errors.ErrKey, err))
		logx.WithContext(ctx).Errorw("查询角色列表失败", fields...)
		return 0, nil, err
	}
	return count, ms, err
}

func (s *RoleService) RoleModelToSub(m models.RoleModel) string {
	return fmt.Sprintf("role_%d", m.Id)
}

func (s *RoleService) AddGroupPolicy(
	ctx context.Context,
	m models.RoleModel,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	sub := s.RoleModelToSub(m)

	// 批量处理权限
	for _, o := range m.Permissions {
		obj := permissionModelToSub(o)
		if err := s.cache.AddGroupPolicy(sub, obj); err != nil {
			logx.WithContext(ctx).Errorw(
				"添加角色与权限的关联策略失败",
				logx.Field(auth.SubKey, sub),
				logx.Field(auth.ObjKey, obj),
				logx.Field("role_id", m.Id),
				logx.Field("permission_id", o.Id),
				logx.Field(errors.ErrKey, err),
			)
			return err
		}
	}

	// 批量处理菜单
	for _, o := range m.Menus {
		obj := menuModelToSub(o)
		if err := s.cache.AddGroupPolicy(sub, obj); err != nil {
			logx.WithContext(ctx).Errorw(
				"添加角色与菜单的关联策略失败",
				logx.Field(auth.SubKey, sub),
				logx.Field(auth.ObjKey, obj),
				logx.Field("role_id", m.Id),
				logx.Field("menu_id", o.Id),
				logx.Field(errors.ErrKey, err),
			)
			return err
		}
	}

	// 批量处理按钮
	for _, o := range m.Buttons {
		obj := buttonModelToSub(o)
		if err := s.cache.AddGroupPolicy(sub, obj); err != nil {
			logx.WithContext(ctx).Errorw(
				"添加角色与按钮的关联策略失败",
				logx.Field(auth.SubKey, sub),
				logx.Field(auth.ObjKey, obj),
				logx.Field("role_id", m.Id),
				logx.Field("button_id", o.Id),
				logx.Field(errors.ErrKey, err),
			)
			return err
		}
	}
	return nil
}

func (s *RoleService) RemoveGroupPolicy(
	ctx context.Context,
	m models.RoleModel,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	sub := s.RoleModelToSub(m)
	// 删除角色作为子级的策略（从其他菜单或权限继承）
	if err := s.cache.RemoveGroupPolicy(0, sub); err != nil {
		logx.WithContext(ctx).Errorw(
			"删除角色作为子级策略失败(该策略继承自其他策略)",
			logx.Field(auth.ObjKey, sub),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}
