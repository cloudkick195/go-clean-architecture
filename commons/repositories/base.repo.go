package repositories

import "gorm.io/gorm"

type IBase interface {
}
type base struct {
	db    *gorm.DB
	table string
}
