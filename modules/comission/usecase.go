package comission

import (
	"context"
	"go_clean_architecture/commons/models"
	"go_clean_architecture/commons/uows"

	"gorm.io/gorm"
)

type IUsecase interface {
	Post(ctx context.Context, tx *gorm.DB, input *CreateInput) (*models.Comission, error)
}

type Usecase struct {
	uow *uows.RepoUOW
}

func NewUsecase(uow *uows.RepoUOW) IUsecase {
	return &Usecase{
		uow: uow,
	}
}

func (u *Usecase) Post(ctx context.Context, tx *gorm.DB, input *CreateInput) (*models.Comission, error) {
	mResult := &models.Comission{
		MemberId:  input.MemberId,
		DepositId: input.DepositId,
		Amount:    input.Amount,
	}
	if err := u.uow.Comission.Create(ctx, tx, mResult); err != nil {
		return nil, err
	}
	return mResult, nil
}
