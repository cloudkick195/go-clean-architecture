package models

import (
	"gorm.io/gorm"
)

type TransferLog struct {
	gorm.Model
	TransferId     *uint
	Log            string
	BankLog        string  `gorm:"type:text"`
	ReqProviderLog *string `gorm:"type:text"`
	ProviderLog    *string `gorm:"type:text"`
	ErrorLog       *string `gorm:"type:text"`
	Ip             string
}

func (TransferLog) TableName() string {
	return "transferLog"
}
