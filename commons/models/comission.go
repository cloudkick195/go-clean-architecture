package models

import "gorm.io/gorm"

type Comission struct {
	gorm.Model
	DepositId uint  `gorm:"column:deposit_id"`
	MemberId  uint  `gorm:"column:member_id"`
	Amount    int64 `gorm:"column:amount;default:0"`
}

func (Comission) TableName() string {
	return "comission"
}
