package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID       uint64    `gorm:"primary_key;column:id;AUTO_INCREMENT;not null" json:"id"`
	CreateAt time.Time `gorm:"autoCreateTime" json:"create_at"`
	UpdateAt time.Time `gorm:"autoUpdateTime" json:"update_at"`
	IsDel    bool      `gorm:"column:is_del;default:0" json:"is_del"`
}

func (b BaseModel) BeforeCreate(tx *gorm.DB) error {
	return nil
}
