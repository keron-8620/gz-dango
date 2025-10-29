package models

import "gz-dango/pkg/database"

type UserModel struct {
	database.StandardModel
	Username string    `gorm:"column:username;type:varchar(50);not null;uniqueIndex;comment:用户名" json:"username"`
	Password string    `gorm:"column:password;type:varchar(150);not null;comment:密码" json:"password"`
	IsActive bool      `gorm:"column:is_active;type:boolean;comment:是否激活" json:"is_active"`
	IsStaff  bool      `gorm:"column:is_staff;type:boolean;comment:是否是工作人员" json:"is_staff"`
	RoleId   uint32    `gorm:"column:role_id;foreignKey:RoleId;references:Id;not null;constraint:OnDelete:CASCADE;comment:角色" json:"role"`
	Role     RoleModel `gorm:"foreignKey:RoleId;constraint:OnDelete:CASCADE"`
}

func (m *UserModel) TableName() string {
	return "customer_user"
}
