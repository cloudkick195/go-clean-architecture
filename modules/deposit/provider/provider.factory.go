package depositProvider

import (
	"errors"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"
)

type IProvider interface {
	Prepare(input *models.Member)
	GetComissionRefer() int64
	Deposit() (string, interface{}, error)
}

type ResponseProvider struct {
	IsSuccess bool
	Data      interface{}
}

type ProviderFactory struct{}

func (factory *ProviderFactory) CreateProvider(config *models.Config, depositModel *models.Deposit) (IProvider, error) {
	if config == nil {
		return nil, commons.ErrInternal(errors.New("invalid provider"))
	}
	switch config.Code {
	case "suplive":
		return NewSupliveProvider(config, depositModel), nil
	default:
		return nil, commons.ErrInternal(errors.New("invalid provider code"))
	}
}
