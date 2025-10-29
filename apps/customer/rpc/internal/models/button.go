package models

import "gz-dango/pkg/database"

type ButtonModel struct {
	database.StandardModel
	Name         string            `gorm:"column:name;type:varchar(50);not null;uniqueIndex;comment:名称" json:"name"`
	ArrangeOrder uint32              `gorm:"column:arrange_order;type:integer;comment:排序" json:"arrange_order"`
	IsActive     bool              `gorm:"column:is_active;type:boolean;comment:是否激活" json:"is_active"`
	Descr        string            `gorm:"column:descr;type:varchar(254);comment:描述" json:"descr"`
	MenuId       uint32              `gorm:"column:menu_id;foreignKey:MenuId;references:Id;not null;constraint:OnDelete:CASCADE;comment:菜单" json:"menu"`
	Menu         MenuModel         `gorm:"foreignKey:MenuId;constraint:OnDelete:CASCADE"`
	Permissions  []PermissionModel `gorm:"many2many:customer_button_permission;joinForeignKey:button_id;joinReferences:permission_id;constraint:OnDelete:CASCADE"`
}

func (m *ButtonModel) TableName() string {
	return "customer_button"
}
