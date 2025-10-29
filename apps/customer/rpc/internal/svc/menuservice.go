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

type MenuService struct {
	gormDB *gorm.DB
	cache  *auth.AuthEnforcer
}

func NewMenuService(
	gormDB *gorm.DB,
	cache *auth.AuthEnforcer,
) *MenuService {
	return &MenuService{
		gormDB: gormDB,
		cache:  cache,
	}
}

func (s *MenuService) CreateModel(ctx context.Context, m *models.MenuModel) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := database.DBCreate(ctx, s.gormDB, &models.MenuModel{}, m); err != nil {
		parent_id := 0
		if m.ParentId != nil {
			parent_id = int(*m.ParentId)
		}
		logx.WithContext(ctx).Errorw(
			"新增菜单模型失败",
			logx.Field("id", m.Id),
			logx.Field("path", m.Path),
			logx.Field("component", m.Component),
			logx.Field("meta", m.Meta.Json()),
			logx.Field("name", m.Name),
			logx.Field("label", m.Label),
			logx.Field("arrange_order", m.ArrangeOrder),
			logx.Field("is_active", m.IsActive),
			logx.Field("descr", m.Descr),
			logx.Field("parent_id", parent_id),
		)
		return err
	}
	return nil
}

func (s *MenuService) UpdateModel(ctx context.Context, data map[string]any, upmap map[string]any, conds ...any) error {
	if err := database.DBUpdate(ctx, s.gormDB, &models.MenuModel{}, data, upmap, conds...); err != nil {
		fields := mapToLogFields(data)
		fields = append(fields, logx.Field(errors.ErrKey, err))
		logx.WithContext(ctx).Errorw("更新菜单模型失败", fields...)
		return err
	}
	return nil
}

func (s *MenuService) DeleteModel(ctx context.Context, conds ...any) error {
	if err := database.DBDelete(ctx, s.gormDB, &models.MenuModel{}, conds...); err != nil {
		logx.WithContext(ctx).Errorw(
			"删除菜单模型失败",
			logx.Field(database.CondsKey, conds),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}

func (s *MenuService) FindModel(
	ctx context.Context,
	preloads []string,
	conds ...any,
) (*models.MenuModel, error) {
	var m models.MenuModel
	if err := database.DBFind(ctx, s.gormDB, preloads, &m, conds...); err != nil {
		logx.WithContext(ctx).Errorw(
			"查询菜单模型失败",
			logx.Field(database.CondsKey, conds),
			logx.Field(errors.ErrKey, err),
		)
		return nil, err

	}
	return &m, nil
}

func (s *MenuService) ListModel(
	ctx context.Context,
	qp database.QueryParams,
) (int64, []models.MenuModel, error) {
	var ms []models.MenuModel
	count, err := database.DBList(ctx, s.gormDB, &models.MenuModel{}, &ms, qp)
	if err != nil {
		fields := qpToLogFields(qp)
		fields = append(fields, logx.Field(errors.ErrKey, err))
		logx.WithContext(ctx).Errorw("查询菜单列表失败", fields...)
		return 0, nil, err
	}
	return count, ms, err
}

func (s *MenuService) ListModelByIds(
	ctx context.Context,
	ids []uint32,
) ([]models.MenuModel, error) {
	if len(ids) == 0 {
		return []models.MenuModel{}, nil
	}
	qp := database.NewPksQueryParams(ids)
	_, ms, err := s.ListModel(ctx, qp)
	return ms, err
}

func (s *MenuService) AddGroupPolicy(
	ctx context.Context,
	m models.MenuModel,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	sub := menuModelToSub(m)

	// 处理父级关系
	if m.Parent != nil {
		obj := menuModelToSub(*m.Parent)
		if err := s.cache.AddGroupPolicy(sub, obj); err != nil {
			logx.WithContext(ctx).Errorw(
				"添加菜单与父级菜单的继承关系策略失败",
				logx.Field(auth.SubKey, sub),
				logx.Field(auth.ObjKey, obj),
				logx.Field("menu_id", m.Id),
				logx.Field("parent_menu_id", m.Parent.Id),
				logx.Field(errors.ErrKey, err),
			)
			return err
		}
	}

	// 批量处理权限
	for _, o := range m.Permissions {
		obj := permissionModelToSub(o)
		if err := s.cache.AddGroupPolicy(sub, obj); err != nil {
			logx.WithContext(ctx).Errorw(
				"添加菜单与权限的关联策略失败",
				logx.Field(auth.SubKey, sub),
				logx.Field(auth.ObjKey, obj),
				logx.Field("menu_id", m.Id),
				logx.Field("permission_id", o.Id),
				logx.Field(errors.ErrKey, err),
			)
			return err
		}
	}
	return nil
}

func (s *MenuService) RemoveGroupPolicy(
	ctx context.Context,
	m models.MenuModel,
	removeInherited bool,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	sub := menuModelToSub(m)
	// 删除该菜单作为子级的策略 (g 规则中的第一个参数)
	if err := s.cache.RemoveGroupPolicy(0, sub); err != nil {
		logx.WithContext(ctx).Errorw(
			"删除菜单作为子级策略失败(该策略继承自其他策略)",
			logx.Field(auth.SubKey, sub),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	if removeInherited {
		// 删除该菜单作为父级的策略 (g 规则中的第二个参数)
		if err := s.cache.RemoveGroupPolicy(1, sub); err != nil {
			logx.WithContext(ctx).Errorw(
				"删除菜单作为父级策略失败(该策略被其他策略继承)",
				logx.Field(auth.ObjKey, sub),
				logx.Field(errors.ErrKey, err),
			)
			return err
		}
	}
	return nil
}

func menuModelToSub(m models.MenuModel) string {
	return fmt.Sprintf("menu_%d", m.Id)
}
