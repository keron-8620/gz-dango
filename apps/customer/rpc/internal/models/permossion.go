package models

import (
	"go-dango/pkg/database"
)

type PermissionModel struct {
	database.StandardModel
	Url    string `gorm:"column:url;type:varchar(150);index:idx_member;comment:URL地址" json:"url"`
	Method string `gorm:"column:method;type:varchar(10);index:idx_member;comment:请求方法" json:"method"`
	Label  string `gorm:"column:label;type:varchar(50);index:idx_member;comment:标签" json:"label"`
	Descr  string `gorm:"column:descr;type:varchar(254);comment:描述" json:"descr"`
}

func (m *PermissionModel) TableName() string {
	return "customer_permission"
}
