package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config là struct chứa các giá trị cấu hình
type Config struct {
	PORT             string `mapstructure:"PORT"`
	DB_CONNECT       string `mapstructure:"DB_CONNECT"`
	DB_POOL_MAX_OPEN int    `mapstructure:"DB_POOL_MAX_OPEN"`
	DB_POOL_MAX_IDLE int    `mapstructure:"DB_POOL_MAX_IDLE"`

	REDIS_ADDRESS           string `mapstructure:"REDIS_ADDRESS"`
	TIME_LIFE_CATCHING      int    `mapstructure:"TIME_LIFE_CATCHING"`
	REDIS_TIME_LIFE_CACHING int    `mapstructure:"REDIS_TIME_LIFE_CACHING"`
	REDIS_PORT              int    `mapstructure:"REDIS_PORT"`
	REDIS_PASSWORD          string `mapstructure:"REDIS_PASSWORD"`

	MAIL_USER     string `mapstructure:"MAIL_USER"`
	MAIL_PASSWORD string `mapstructure:"MAIL_PASSWORD"`
	MAIL_HOST     string `mapstructure:"MAIL_HOST"`
	MAIL_PORT     string `mapstructure:"MAIL_PORT"`

	WEBSITE_URL string `mapstructure:"WEBSITE_URL"`

	AUTH_TOKEN      string `mapstructure:"AUTH_TOKEN"`
	BANK_AUTH_TOKEN string `mapstructure:"BANK_AUTH_TOKEN"`
	AUTHORIZED_IPS  string `mapstructure:"AUTHORIZED_IPS"`
}

var Env Config

// LoadConfig là hàm dùng để load file cấu hình và lưu các giá trị vào Config
func InitConfig() {
	var config Config

	viper.SetConfigFile(".env") // Đặt tên file cấu hình là .env
	viper.AllowEmptyEnv(true)   // Cho phép biến môi trường rỗng
	viper.AutomaticEnv()        // Tự động đọc biến môi trường

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}

	// Map các giá trị từ file config vào struct Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %s ", err))
	}
	Env = config
}
