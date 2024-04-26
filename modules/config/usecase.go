package configModule

import (
	"context"
	"go_clean_architecture/commons/models"
	"go_clean_architecture/commons/uows"
)

type IUsecase interface {
	Get(ctx context.Context, code string) (*models.Config, error)
	List(ctx context.Context) (list []models.Config, err error)
}

func NewUsecase(uow *uows.RepoUOW) IUsecase {
	return &Usecase{
		uow: uow,
	}
}

type Usecase struct {
	uow *uows.RepoUOW
}

func (a *Usecase) Get(ctx context.Context, code string) (*models.Config, error) {
	return a.uow.Config.Get(ctx, &models.Config{Code: code})
}
func (a *Usecase) List(ctx context.Context) (list []models.Config, err error) {
	return a.uow.Config.List(ctx)
}
