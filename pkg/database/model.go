package database

import (
	"time"
)

type BaseModel struct {
	Id uint32 `gorm:"column:id;primary_key;AUTO_INCREMENT;comment:编号" json:"id"`
}

type StandardModel struct {
	BaseModel

	CreatedAt time.Time `gorm:"column:created_at;comment:创建时间;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;comment:修改时间;autoUpdateTime" json:"updated_at"`
}
