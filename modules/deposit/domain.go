package deposit

import (
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"
	depositProvider "go_clean_architecture/modules/deposit/provider"

	"time"
)

type CreateInput struct {
	SupId         int64 `validate:"required,min=5"`
	Num           int64 `validate:"required,gt=0"`
	ComissionType *int8 `validate:"omitempty,oneof=0 1"`
	Member        *models.Member
}

type DepositInput struct {
	Provider       depositProvider.IProvider
	Deposit        *models.Deposit
	ComissionRefer int64
	AmountStrategy DepositStrategy
}

type Output struct {
	Id          uint
	Code        string
	ProviderID  string
	Amount      int64
	CoinAmount  int64
	PriceCoin   int64
	CompletedAt *time.Time
	CreatedAt   *time.Time
}
type Query struct {
	commons.Pagination
	models.FilterDeposit
}
type QueryInput struct {
	CommisionType *int8
	MemberId      uint
	From          string
	To            string
}
