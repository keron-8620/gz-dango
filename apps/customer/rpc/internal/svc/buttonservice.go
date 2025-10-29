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

type ButtonService struct {
	gormDB *gorm.DB
	cache  *auth.AuthEnforcer
}

func NewButtonService(
	gormDB *gorm.DB,
	cache *auth.AuthEnforcer,
) *ButtonService {
	return &ButtonService{
		gormDB: gormDB,
		cache:  cache,
	}
}

func (s *ButtonService) CreateModel(ctx context.Context, m *models.ButtonModel) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := database.DBCreate(ctx, s.gormDB, &models.ButtonModel{}, m); err != nil {
		logx.WithContext(ctx).Errorw(
			"新增按钮模型失败",
			logx.Field("name", m.Name),
			logx.Field("arrange_order", m.ArrangeOrder),
			logx.Field("is_active", m.IsActive),
			logx.Field("descr", m.Descr),
			logx.Field("menu_id", m.MenuId),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}

func (s *ButtonService) UpdateModel(ctx context.Context, data map[string]any, upmap map[string]any, conds ...any) error {
	if err := database.DBUpdate(ctx, s.gormDB, &models.ButtonModel{}, data, upmap, conds...); err != nil {
		fields := mapToLogFields(data)
		fields = append(fields, logx.Field(errors.ErrKey, err))
		logx.WithContext(ctx).Errorw("更新按钮模型失败", fields...)
	}
	return nil
}

func (s *ButtonService) DeleteModel(ctx context.Context, conds ...any) error {
	if err := database.DBDelete(ctx, s.gormDB, &models.ButtonModel{}, conds...); err != nil {
		logx.WithContext(ctx).Errorw(
			"删除按钮模型失败",
			logx.Field(database.CondsKey, conds),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	return nil
}

func (s *ButtonService) FindModel(
	ctx context.Context,
	preloads []string,
	conds ...any,
) (*models.ButtonModel, error) {
	var m models.ButtonModel
	if err := database.DBFind(ctx, s.gormDB, preloads, &m, conds...); err != nil {
		logx.WithContext(ctx).Errorw(
			"查询按钮模型失败",
			logx.Field(database.CondsKey, conds),
			logx.Field(errors.ErrKey, err),
		)
		return nil, err
	}
	return &m, nil
}

func (s *ButtonService) ListModel(
	ctx context.Context,
	qp database.QueryParams,
) (int64, []models.ButtonModel, error) {
	var ms []models.ButtonModel
	count, err := database.DBList(ctx, s.gormDB, &models.ButtonModel{}, &ms, qp)
	if err != nil {
		fields := qpToLogFields(qp)
		fields = append(fields, logx.Field(errors.ErrKey, err))
		logx.WithContext(ctx).Errorw("查询按钮列表失败", fields...)
		return 0, nil, err
	}
	return count, ms, err
}

func (s *ButtonService) ListModelByIds(
	ctx context.Context,
	ids []uint32,
) ([]models.ButtonModel, error) {
	if len(ids) == 0 {
		return []models.ButtonModel{}, nil
	}
	qp := database.NewPksQueryParams(ids)
	_, ms, err := s.ListModel(ctx, qp)
	return ms, err
}

func (s *ButtonService) AddGroupPolicy(
	ctx context.Context,
	m models.ButtonModel,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	sub := buttonModelToSub(m)
	menuObj := menuModelToSub(m.Menu)
	if err := s.cache.AddGroupPolicy(sub, menuObj); err != nil {
		logx.WithContext(ctx).Errorw(
			"添加按钮与菜单的继承关系策略失败",
			logx.Field(auth.SubKey, sub),
			logx.Field(auth.ObjKey, menuObj),
			logx.Field("button_id", m.Id),
			logx.Field("menu_id", m.MenuId),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}

	// 批量处理权限
	for _, o := range m.Permissions {
		obj := permissionModelToSub(o)
		if err := s.cache.AddGroupPolicy(sub, obj); err != nil {
			logx.WithContext(ctx).Errorw(
				"添加按钮与权限的关联策略失败",
				logx.Field(auth.SubKey, sub),
				logx.Field(auth.ObjKey, obj),
				logx.Field("button_id", m.Id),
				logx.Field("permission_id", o.Id),
				logx.Field(errors.ErrKey, err),
			)
			return err
		}
	}
	return nil
}

func (s *ButtonService) RemoveGroupPolicy(
	ctx context.Context,
	m models.ButtonModel,
	removeInherited bool,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	sub := buttonModelToSub(m)
	// 删除该按钮作为子级的策略（从其他菜单或权限继承）
	if err := s.cache.RemoveGroupPolicy(0, sub); err != nil {
		logx.WithContext(ctx).Errorw(
			"删除按钮作为子级策略失败(该策略继承自其他策略)",
			logx.Field(auth.ObjKey, sub),
			logx.Field(errors.ErrKey, err),
		)
		return err
	}
	if removeInherited {
		// 删除该按钮作为父级的策略（被其他菜单或权限继承）
		if err := s.cache.RemoveGroupPolicy(1, sub); err != nil {
			logx.WithContext(ctx).Errorw(
				"删除按钮作为父级策略失败(该策略被其他策略继承)",
				logx.Field(auth.ObjKey, sub),
				logx.Field(errors.ErrKey, err),
			)
			return err
		}
	}
	return nil
}

func buttonModelToSub(m models.ButtonModel) string {
	return fmt.Sprintf("button_%d", m.Id)
}
