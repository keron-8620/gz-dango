package database

import (
	"time"
)

type BaseModel struct {
	Id uint32 `gorm:"column:id;primary_key;AUTO_INCREMENT;comment:编号" json:"id"`
}

type StandardModel struct {
	BaseModel

	CreatedAt time.Time `gorm:"column:create_at;comment:创建时间" json:"create_at"`
	UpdatedAt time.Time `gorm:"column:update_at;comment:修改时间" json:"update_at"`
}
