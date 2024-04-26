package models

import (
	"errors"
	"go_clean_architecture/commons"
	"time"

	"gorm.io/gorm"
)

type Deposit struct {
	gorm.Model
	MemberId      uint `gorm:"member_id"`
	Member        *Member
	ComissionType *int8      `gorm:"is_comission"`
	ConfigCode    string     `gorm:"config_code"`
	ProviderID    string     `gorm:"provider_id"`
	Amount        int64      `gorm:"column:amount;default:0"`
	CoinAmount    int64      `gorm:"column:coin_amount;default:0"`
	TotalCoin     int64      `gorm:"column:total_coin;default:0"`
	PriceCoin     int64      `gorm:"column:price_coin;default:0"`
	Commission    int64      `gorm:"column:commission;default:0"`
	CompletedAt   *time.Time `gorm:"column:completed_at"`
}

type FilterDeposit struct {
	ComissionType *int8 `json:"ComissionType" form:"ComissionType"`
	MemberId      *uint `json:"MemberId" form:"MemberId"`
	From          string
	To            string
}

func (Deposit) TableName() string {
	return "deposit"
}

var (
	ErrWrongDepositID = commons.NewCustomError(
		errors.New("ErrWrongDepositID"),
		"ErrWrongDepositID",
		"ErrWrongDepositID",
	)
)
