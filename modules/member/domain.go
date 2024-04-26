package member

import (
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"
	"go_clean_architecture/modules/deposit"
	"go_clean_architecture/modules/transfer"
)

type CreateInput struct {
	SupId        int64  `validate:"required,min=5"`
	FullName     string `validate:"required"`
	Email        string `validate:"required,email"`
	Password     string `validate:"required,min=5"`
	Ip           string `validate:"required,ip"`
	ReferralCode *string
}

type DepositInput struct {
	SupId         int64 `validate:"required,min=5"`
	Num           int64 `validate:"required,gt=0"`
	ComissionType *int8 `validate:"omitempty,oneof=1"`
}

type Ouput struct {
	SupId            int64
	FullName         string
	Email            string
	Amount           int64
	Commission       int64
	TotalAmount      int64
	TotalTransaction int64
	Role             string
	Status           models.MemberStatusEnum
	IsAgency         bool
}

type UpdatePasswordInput struct {
	Id          uint
	Password    string `validate:"required,min=5"`
	NewPassword string `validate:"required,min=5"`
}

type UpdateMember struct {
	Id       uint
	Status   string `validate:"required,oneof=Active Blocked"`
	IsAgency bool
}

type Query struct {
	commons.Pagination
	models.FilterMember
}

type TransactionHistoryQuery struct {
	Id   uint
	From string `form:"from"`
	To   string `form:"to"`
}

type TransactionHistoryOutput struct {
	Transfer []transfer.Output
	Deposit  []deposit.Output
}
