package deposit

import (
	"fmt"
	"go_clean_architecture/commons/models"

	"gorm.io/gorm"
)

type DepositStrategy interface {
	GetAmount(input *CreateInput) int64
	SetDepositRecord(input *models.Deposit)
	IsDeposit() bool
	UpdateMember(inputMember *models.Member, inputDeposit *models.Deposit) map[string]interface{}
}

type AmountDepositStrategy struct {
}

func (a *AmountDepositStrategy) GetAmount(input *CreateInput) int64 {
	return input.Member.Amount
}

func (a *AmountDepositStrategy) SetDepositRecord(input *models.Deposit) {

}

func (a *AmountDepositStrategy) IsDeposit() bool {
	return true
}

func (a *AmountDepositStrategy) UpdateMember(inputMember *models.Member, inputDeposit *models.Deposit) map[string]interface{} {
	return map[string]interface{}{
		"Amount":                 gorm.Expr(fmt.Sprintf("amount - %d", inputDeposit.Amount)),
		"Commission":             gorm.Expr(fmt.Sprintf("commission + %d", inputDeposit.Commission)),
		"TotalTransactionAmount": gorm.Expr(fmt.Sprintf("total_transaction_amount + %d", inputDeposit.Amount)),
		"TotalTransaction":       gorm.Expr(fmt.Sprintf("total_transaction + %d", 1)),
	}
}

type CommissionDepositStrategy struct {
}

func (c *CommissionDepositStrategy) GetAmount(input *CreateInput) int64 {
	return input.Member.Commission
}
func (a *CommissionDepositStrategy) SetDepositRecord(input *models.Deposit) {
	input.Commission = 0
	input.CompletedAt = nil
}
func (a *CommissionDepositStrategy) IsDeposit() bool {
	return false
}
func (a *CommissionDepositStrategy) UpdateMember(inputMember *models.Member, inputDeposit *models.Deposit) map[string]interface{} {
	return map[string]interface{}{
		"Commission":             gorm.Expr(fmt.Sprintf("commission - %d", inputDeposit.Amount)),
		"TotalTransactionAmount": gorm.Expr(fmt.Sprintf("total_transaction_amount + %d", inputDeposit.Amount)),
		"TotalTransaction":       gorm.Expr(fmt.Sprintf("total_transaction + %d", 1)),
	}
}

type DepositStrategyFactory struct{}

func (f *DepositStrategyFactory) CreateStrategy(ComissionType *int8) DepositStrategy {
	if ComissionType != nil {
		if *ComissionType == 1 {
			return &CommissionDepositStrategy{}
		}
	}
	return &AmountDepositStrategy{}
}
