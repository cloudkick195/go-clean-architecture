package auth

type LoginInput struct {
	SupId    int64  `validate:"required,min=5"`
	Password string `validate:"required,min=5"`
}
