package comission

type CreateInput struct {
	DepositId uint  `gorm:"column:deposit_id"`
	MemberId  uint  `gorm:"column:member_id"`
	Amount    int64 `gorm:"column:amount"`
}
