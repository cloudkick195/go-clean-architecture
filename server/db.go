package server

import (
	"fmt"
	"go_clean_architecture/commons/models"
	"go_clean_architecture/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Ghi log ra stdout
		logger.Config{
			SlowThreshold: time.Second, // Thời gian mà các truy vấn được xem như chậm (nếu thực thi lâu hơn, sẽ được ghi lại trong log)
			LogLevel:      logger.Info, // Mức độ log
			Colorful:      true,        // Sử dụng màu sắc cho log
		},
	)

	dbURI := fmt.Sprintf(
		"%s",
		config.Env.DB_CONNECT,
	)

	conn, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{Logger: newLogger})

	if err != nil {
		log.Fatalln(err)
	}
	DbMigrate(conn)
	SeedData(conn)
	return conn
}

func DbMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Log{},
		&models.Member{},
		&models.Deposit{},
		&models.Comission{},
		&models.DepositLog{},
		&models.Transfer{},
		&models.TransferLog{},
		&models.Config{},
	)
	if err != nil {
		panic("Migration error")
	}
}

func SeedData(db *gorm.DB) {

}
