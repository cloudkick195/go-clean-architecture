package log

type ReqQuery struct {
	Status *int `query:"status" validate:"omitempty,oneof=0 1"`
}

type ReqParams struct {
	Id int
}

// CreateIncomeExpenseCategoryRequest represents the request payload for creating a new income/expense category
type CreateInput struct {
	StatusCode int
	RootErr    string
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"key"`
	Api        string `json:"api"`
	Request    string `json:"request"`
	Ip         string `json:"ip"`
}

// UpdateIncomeExpenseCategoryRequest represents the request payload for updating an existing income/expense category
type UpdateInput struct {
	Id      uint
	RootErr string `gorm:"type:text"`
	Message string `json:"message"`
	Log     string `json:"log"`
	Key     string `json:"key"`
}

type Output struct {
	StatusCode int
	Message    string `json:"message"`
	Key        string `json:"key"`
}
