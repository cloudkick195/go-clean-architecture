package models

import (
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	StatusCode int
	RootErr    string `gorm:"type:text"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"key"`
	Api        string `json:"api"`       // Thêm trường để lưu thông tin lỗi API
	Request    string `gorm:"type:text"` // Thêm trường để lưu thông tin request
	Ip         string
}

func (Log) TableName() string {
	return "log"
}
