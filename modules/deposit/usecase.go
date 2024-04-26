package deposit

import (
	"context"
	"errors"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"
	"go_clean_architecture/commons/uows"
	configModule "go_clean_architecture/modules/config"
	depositProvider "go_clean_architecture/modules/deposit/provider"
	"go_clean_architecture/modules/depositLog"
	"go_clean_architecture/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type IUsecase interface {
	Post(ctx context.Context, tx *gorm.DB, input *CreateInput) (*DepositInput, error)
	Deposit(ctx context.Context, input *DepositInput) error
	GetByFromTo(ctx context.Context, input *models.FilterDeposit) (result []Output, err error)
	Detail(ctx context.Context, input uint) (*models.Deposit, error)
	List(ctx context.Context, paging *commons.Pagination, queries *models.FilterDeposit) (results []models.Deposit, err error)
	UpdateDeposit(ctx context.Context, input uint) error
}

type Usecase struct {
	uow                    *uows.RepoUOW
	depositLog             depositLog.IUsecase
	configModule           configModule.IUsecase
	providerFactory        *depositProvider.ProviderFactory
	depositStrategyFactory DepositStrategyFactory
}

func NewUsecase(uow *uows.RepoUOW, depositLog depositLog.IUsecase, configModule configModule.IUsecase, providerFactory *depositProvider.ProviderFactory, depositStrategyFactory DepositStrategyFactory) IUsecase {
	return &Usecase{
		uow:                    uow,
		depositLog:             depositLog,
		configModule:           configModule,
		providerFactory:        providerFactory,
		depositStrategyFactory: depositStrategyFactory,
	}
}

func (u *Usecase) List(ctx context.Context, paging *commons.Pagination, queries *models.FilterDeposit) (results []models.Deposit, err error) {
	paging.HandleRequest()
	return u.uow.Deposit.List(ctx, paging, queries)
}

func (u *Usecase) Detail(ctx context.Context, input uint) (*models.Deposit, error) {
	result, err := u.uow.Deposit.GetById(ctx, input)
	if err != nil {
		return nil, err
	}
	resultMember, _ := u.uow.Member.GetById(ctx, result.MemberId)
	result.Member = resultMember
	return result, nil
}

func (u *Usecase) GetByFromTo(ctx context.Context, input *models.FilterDeposit) (result []Output, err error) {
	deposits, err := u.uow.Deposit.GetByFromTo(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, deposit := range deposits {
		output := Output{
			Id:          deposit.ID,
			Code:        deposit.ConfigCode,
			ProviderID:  deposit.ProviderID,
			Amount:      deposit.Amount,
			CoinAmount:  deposit.CoinAmount,
			CompletedAt: deposit.CompletedAt,
			CreatedAt:   &deposit.CreatedAt,
		}
		result = append(result, output)
	}

	return result, nil
}

func (u *Usecase) Post(ctx context.Context, tx *gorm.DB, input *CreateInput) (*DepositInput, error) {
	now := time.Now()
	if err := utils.Validate.Struct(input); err != nil {
		return nil, commons.ErrInvalidRequest(err)
	}

	amountStrategy := u.depositStrategyFactory.CreateStrategy(input.ComissionType)

	if amountStrategy.GetAmount(input) < input.Num {
		return nil, commons.NewCustomError(
			errors.New("amount must be greater than"),
			"amount must be greater than",
			"ErrAmount",
		)
	}

	configResult, err := u.configModule.Get(ctx, "suplive")
	if err != nil {
		return nil, err
	}

	if configResult == nil {
		return nil, commons.ErrDB(errors.New("config not found"))
	}

	mResult := &models.Deposit{
		MemberId:      input.Member.ID,
		ProviderID:    strconv.FormatInt(input.SupId, 10),
		ComissionType: input.ComissionType,
		Amount:        input.Num,
		CompletedAt:   &now,
	}

	providerStruct, err := u.providerFactory.CreateProvider(configResult, mResult)
	if err != nil {
		return nil, err
	}
	if providerStruct != nil {
		mResult.ConfigCode = configResult.Code
		providerStruct.Prepare(input.Member)
		amountStrategy.SetDepositRecord(mResult)
	}

	if err := u.uow.Deposit.Create(ctx, tx, mResult); err != nil {
		if err := u.depositLog.CreateWithParams(ctx, nil, err, mResult); err != nil {
			return nil, err
		}
		return nil, err
	}

	return &DepositInput{
		Provider:       providerStruct,
		Deposit:        mResult,
		ComissionRefer: providerStruct.GetComissionRefer(),
		AmountStrategy: amountStrategy,
	}, nil
}

func (u *Usecase) Deposit(ctx context.Context, input *DepositInput) error {
	reqProvider, resProvider, err := input.Provider.Deposit()

	if err != nil {
		if err := u.depositLog.CreateWithParams(ctx, nil, err, reqProvider); err != nil {
			return err
		}
		return err
	}

	if err := u.depositLog.CreateWithParams(ctx, &input.Deposit.ID, resProvider, reqProvider); err != nil {
		return err
	}
	return nil
}

func (u *Usecase) UpdateDeposit(ctx context.Context, input uint) error {
	result, err := u.uow.Deposit.GetById(ctx, input)
	if err != nil {
		return err
	}
	if result == nil {
		return commons.ErrDB(errors.New("deposit not found"))
	}
	if result.CompletedAt != nil {
		return DepositHasBeenCompleted
	}

	now := time.Now()
	return u.uow.Do(ctx, func(tx *gorm.DB) error {
		if err := u.uow.Deposit.Update(ctx, tx, result, map[string]interface{}{"CompletedAt": &now}); err != nil {
			return err
		}
		configResult, err := u.configModule.Get(ctx, result.ConfigCode)
		if err != nil {
			return err
		}

		if configResult == nil {
			return commons.ErrDB(errors.New("config not found"))
		}
		providerStruct, err := u.providerFactory.CreateProvider(configResult, result)
		if err != nil {
			return err
		}
		memberDeposit := &DepositInput{
			Provider: providerStruct,
			Deposit:  result,
		}

		if err := u.Deposit(ctx, memberDeposit); err != nil {
			return err
		}
		return nil
	})
}
