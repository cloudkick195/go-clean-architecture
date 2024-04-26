package models

import (
	"gorm.io/gorm"
)

type DepositLog struct {
	gorm.Model
	DepositId *uint
	Log       string `gorm:"type:text"`
	Req       string `gorm:"type:text"`
}

func (DepositLog) TableName() string {
	return "depositLog"
}
