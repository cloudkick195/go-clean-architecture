package depositLog

import (
	"context"
	"go_clean_architecture/commons/models"
	"go_clean_architecture/commons/uows"
	"go_clean_architecture/utils"
)

type IUsecase interface {
	Get(ctx context.Context, id uint) (*models.DepositLog, error)
	Create(ctx context.Context, input *CreateInput) (*models.DepositLog, error)
	CreateWithParams(ctx context.Context, depositId *uint, log interface{}, req interface{}) error
}

func NewUsecase(uow *uows.RepoUOW) IUsecase {
	return &Usecase{
		uow: uow,
	}
}

type Usecase struct {
	uow *uows.RepoUOW
}

func (a *Usecase) Get(ctx context.Context, id uint) (*models.DepositLog, error) {
	return a.uow.DepositLog.GetById(ctx, id)
}

func (a *Usecase) CreateWithParams(ctx context.Context, depositId *uint, log interface{}, req interface{}) error {
	inputdepositLog := &CreateInput{
		DepositId: depositId,
		Log:       utils.InterfaceToString(log),
		Req:       utils.InterfaceToString(req),
	}

	if _, err := a.Create(ctx, inputdepositLog); err != nil {
		return err
	}
	return nil
}

func (a *Usecase) Create(ctx context.Context, input *CreateInput) (*models.DepositLog, error) {
	m := &models.DepositLog{
		DepositId: input.DepositId,
		Log:       input.Log,
		Req:       input.Req,
	}
	if err := a.uow.DepositLog.Create(ctx, m); err != nil {
		return nil, err
	}

	return m, nil
}
