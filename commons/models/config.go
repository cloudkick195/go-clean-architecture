package models

type Config struct {
	Id       int `gorm:"column:id;primaryKey"`
	Code     string
	Name     string
	Config   string `gorm:"type:text"`
	IsActive bool
}

func (Config) TableName() string {
	return "config"
}
