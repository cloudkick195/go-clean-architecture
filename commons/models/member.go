package models

import (
	"go_clean_architecture/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MemberStatusEnum string

const (
	MemberStatusActive   MemberStatusEnum = "Active"
	MemberStatusInActive MemberStatusEnum = "InActive"
	MemberStatusBlocked  MemberStatusEnum = "Blocked"
)

type Member struct {
	gorm.Model
	SupId                  int64  `json:"sup_id" gorm:"column:sup_id"`
	FullName               string `json:"full_name" gorm:"column:full_name"`
	Email                  string `json:"email" gorm:"column:email"`
	Password               string `json:"-" gorm:"column:password"`
	Salt                   string `json:"-" gorm:"column:salt"`
	Token                  string `json:"-" gorm:"column:token"`
	TokenExpiry            int64
	Amount                 int64 `json:"amount" gorm:"column:amount;default:0"`
	Commission             int64 `gorm:"column:commission;default:0"`
	TotalTransactionAmount int64 `gorm:"column:total_transaction_amount;default:0"`
	TotalTransaction       int64 `gorm:"column:total_transaction;default:0"`
	IsAgency               bool
	Status                 MemberStatusEnum `json:"status" gorm:"column:status"`
	Ip                     string           `json:"ip" gorm:"column:ip"`
	ReferralCode           *string          `json:"referral_code" gorm:"column:referral_code"`
	Role                   string           `json:"role" gorm:"column:role"`
}

type FilterMember struct {
	SupId *string
}

// TableName sets the table name for the GORM model.
func (m *Member) TableName() string {
	return "member"
}

func (m *Member) BeforeCreate(tx *gorm.DB) error {
	m.IsAgency = false
	return nil
}

func (u *Member) HashPassword() error {
	salt := utils.GenSalt(20)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Salt = salt
	u.Password = string(hashedPassword)
	return nil
}
func (u *Member) GetPassword() string {
	return u.Password
}

func (u *Member) ComparePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password+u.Salt)); err != nil {
		return false
	}
	return true
}
