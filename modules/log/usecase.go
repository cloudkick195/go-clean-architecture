package log

import (
	"context"
	"encoding/json"
	"go_clean_architecture/commons/models"
	"go_clean_architecture/commons/uows"
)

type IUsecase interface {
	Get(ctx context.Context, id uint) (*models.Log, error)
	Create(ctx context.Context, input *CreateInput) error
	Update(ctx context.Context, input *UpdateInput) error
}

func NewUsecase(uow *uows.RepoUOW) IUsecase {
	return &Usecase{
		uow: uow,
	}
}

type Usecase struct {
	uow *uows.RepoUOW
}

func (a *Usecase) Get(ctx context.Context, id uint) (*models.Log, error) {
	return a.uow.Log.GetById(ctx, id)
}

func (a *Usecase) Create(ctx context.Context, input *CreateInput) error {

	m := &models.Log{
		StatusCode: input.StatusCode,
		RootErr:    input.RootErr,
		Message:    input.Message,
		Log:        input.Log,
		Key:        input.Key,
		Api:        input.Api,
		Request:    input.Request,
		Ip:         input.Ip,
	}
	return a.uow.Log.Create(ctx, m)
}

func (a *Usecase) Update(ctx context.Context, input *UpdateInput) error {
	schema, err := a.uow.Log.GetById(ctx, input.Id)
	if err != nil {
		return err
	}
	var mapping map[string]interface{}
	jsonInput, _ := json.Marshal(input)
	json.Unmarshal(jsonInput, &mapping)
	return a.uow.Log.Updates(ctx, schema, input.Id, mapping)
}
