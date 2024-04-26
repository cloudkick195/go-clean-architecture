package models

import (
	"time"

	"gorm.io/gorm"
)

type TransferStatusEnum string

var TransferStatus = map[TransferStatusEnum]TransferStatusEnum{
	"Pending": TransferStatusEnum("Pending"),
	"Success": TransferStatusEnum("Success"),
	"Fail":    TransferStatusEnum("Fail"),
}

type Transfer struct {
	gorm.Model
	PrivateID             int64              `gorm:"primary_key"`
	ConfigCode            string             `gorm:"config_code"`
	ProviderID            string             `gorm:"provider_id"`
	BankID                string             `gorm:"column:tid"`
	Description           string             `gorm:"column:description"`
	Amount                float64            `gorm:"column:amount"`
	CoinAmount            float64            `gorm:"column:coint_amount"`
	PriceCoin             float64            `gorm:"column:price_coin"`
	Commission            float64            `gorm:"column:commission"`
	TotalCoin             float64            `gorm:"column:total_coin"`
	CusumBalance          float64            `gorm:"column:cusum_balance"`
	When                  time.Time          `gorm:"column:when"`
	BankSubAccID          string             `gorm:"column:bank_sub_acc_id"`
	SubAccID              string             `gorm:"column:sub_acc_id"`
	BankName              string             `gorm:"column:bank_name"`
	BankAbbreviation      string             `gorm:"column:bank_abbreviation"`
	VirtualAccount        string             `gorm:"column:virtual_account"`
	VirtualAccountName    string             `gorm:"column:virtual_account_name"`
	CorresponsiveName     string             `gorm:"column:corresponsive_name"`
	CorresponsiveAccount  string             `gorm:"column:corresponsive_account"`
	CorresponsiveBankID   string             `gorm:"column:corresponsive_bank_id"`
	CorresponsiveBankName string             `gorm:"column:corresponsive_bank_name"`
	TransferStatus        TransferStatusEnum `gorm:"column:transfer_status"`
	CompletedAt           *time.Time         `gorm:"column:completed_at"`
}

type FilterTransfer struct {
	ProviderID string
	From       string
	To         string
}

func (m *Transfer) BeforeCreate(tx *gorm.DB) error {
	m.TransferStatus = TransferStatus["Pending"]
	return nil
}

func (Transfer) TableName() string {
	return "transfer"
}
