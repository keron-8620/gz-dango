package models

import "go-dango/pkg/database"

type RoleModel struct {
	database.StandardModel
	Name        string            `gorm:"column:name;type:varchar(50);not null;uniqueIndex;comment:名称" json:"name"`
	Descr       string            `gorm:"column:descr;type:varchar(254);comment:描述" json:"descr"`
	Permissions []PermissionModel `gorm:"many2many:customer_role_permission;joinForeignKey:role_id;joinReferences:permission_id;constraint:OnDelete:CASCADE"`
	Menus       []MenuModel       `gorm:"many2many:customer_role_menu;joinForeignKey:role_id;joinReferences:menu_id;constraint:OnDelete:CASCADE"`
	Buttons     []ButtonModel     `gorm:"many2many:customer_role_button;joinForeignKey:role_id;joinReferences:button_id;constraint:OnDelete:CASCADE"`
}

func (m *RoleModel) TableName() string {
	return "customer_role"
}
