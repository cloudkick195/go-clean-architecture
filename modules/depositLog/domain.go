package depositLog

type ReqQuery struct {
	Status *int `query:"status" validate:"omitempty,oneof=0 1"`
}

type ReqParams struct {
	Id int
}

// CreateIncomeExpenseCategoryRequest represents the request payload for creating a new income/expense category
type CreateInput struct {
	DepositId *uint
	Log       string `gorm:"type:text"`
	Req       string `gorm:"type:text"`
}
