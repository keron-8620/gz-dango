package models

import (
	"encoding/json"

	"gz-dango/pkg/database"
)

type Meta struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
}

func (m *Meta) Json() string {
	metaBytes, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(metaBytes)
}

type MenuModel struct {
	database.StandardModel
	Path         string            `gorm:"column:path;type:varchar(100);not null;uniqueIndex;comment:前端路由" json:"url"`
	Component    string            `gorm:"column:component;type:varchar(200);not null;comment:请求方式" json:"method"`
	Name         string            `gorm:"column:name;type:varchar(50);not null;uniqueIndex;comment:名称" json:"name"`
	Meta         Meta              `gorm:"column:meta;serializer:json;comment:菜单信息" json:"meta"`
	Label        string            `gorm:"column:label;type:varchar(50);index:idx_member;comment:标签" json:"label"`
	ArrangeOrder uint32            `gorm:"column:arrange_order;type:integer;comment:排序" json:"arrange_order"`
	IsActive     bool              `gorm:"column:is_active;type:boolean;comment:是否激活" json:"is_active"`
	Descr        string            `gorm:"column:descr;type:varchar(254);comment:描述" json:"descr"`
	ParentId     *uint32           `gorm:"column:parent_id;foreignKey:ParentId;references:Id;constraint:OnDelete:CASCADE;comment:菜单" json:"parent"`
	Parent       *MenuModel        `gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE"`
	Permissions  []PermissionModel `gorm:"many2many:customer_menu_permission;joinForeignKey:menu_id;joinReferences:permission_id;constraint:OnDelete:CASCADE"`
}

func (m *MenuModel) TableName() string {
	return "customer_menu"
}
