package member

import (
	"context"
	"fmt"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/middlewares"
	"go_clean_architecture/commons/models"
	"go_clean_architecture/commons/uows"
	"go_clean_architecture/modules/comission"
	"go_clean_architecture/modules/deposit"
	"go_clean_architecture/modules/transfer"
	"go_clean_architecture/utils"
	"strconv"
	"sync"

	"gorm.io/gorm"
)

type IUsecase interface {
	Create(ctx context.Context, input *CreateInput) error
	Get(ctx context.Context, input *models.Member) (*models.Member, error)
	GetById(ctx context.Context, input uint) (*models.Member, error)
	Profile(ctx context.Context, input uint) (*Ouput, error)
	Deposit(ctx context.Context, id uint, input *DepositInput, mw middlewares.IMiddlewareManager) error
	Update(ctx context.Context, tx *gorm.DB, inputWhere *models.Member, inputUpdate map[string]interface{}) error
	UpdatePassword(ctx context.Context, input *UpdatePasswordInput) error
	UpdateMember(ctx context.Context, input *UpdateMember) error
	TransactionHistory(ctx context.Context, input *TransactionHistoryQuery) (*TransactionHistoryOutput, error)
	List(ctx context.Context, paging *commons.Pagination, queries *models.FilterMember) (results []models.Member, err error)
}

type Usecase struct {
	uow       *uows.RepoUOW
	deposit   deposit.IUsecase
	transfer  transfer.IUsecase
	comission comission.IUsecase
}

func NewUsecase(uow *uows.RepoUOW, deposit deposit.IUsecase, transfer transfer.IUsecase, comission comission.IUsecase) IUsecase {
	return &Usecase{
		uow:       uow,
		deposit:   deposit,
		transfer:  transfer,
		comission: comission,
	}
}

func (a *Usecase) GetById(ctx context.Context, input uint) (*models.Member, error) {
	return a.uow.Member.GetById(ctx, input)
}

func (a *Usecase) Get(ctx context.Context, input *models.Member) (*models.Member, error) {
	return a.uow.Member.Get(ctx, input)
}
func (a *Usecase) List(ctx context.Context, paging *commons.Pagination, queries *models.FilterMember) (results []models.Member, err error) {
	paging.HandleRequest()
	return a.uow.Member.List(ctx, paging, queries)
}

func (a *Usecase) Profile(ctx context.Context, input uint) (*Ouput, error) {
	result, err := a.uow.Member.GetById(ctx, input)
	if err != nil {
		return nil, err
	}
	output := &Ouput{
		SupId:            result.SupId,
		FullName:         result.FullName,
		Email:            result.Email,
		Amount:           result.Amount,
		Commission:       result.Commission,
		TotalAmount:      result.TotalTransactionAmount,
		TotalTransaction: result.TotalTransaction,
		Role:             result.Role,
		IsAgency:         result.IsAgency,
		Status:           result.Status,
	}
	return output, nil
}

func (a *Usecase) TransactionHistory(ctx context.Context, input *TransactionHistoryQuery) (*TransactionHistoryOutput, error) {
	resultUser, err := a.GetById(ctx, input.Id)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var resultDeposit []deposit.Output
	var resultTransfer []transfer.Output
	var errVar error
	wg.Add(2)
	go func() {
		result, err := a.deposit.GetByFromTo(ctx, &models.FilterDeposit{
			MemberId: &input.Id,
			From:     input.From,
			To:       input.To,
		})
		if err != nil {
			errVar = err
		}
		resultDeposit = result
		wg.Done()
	}()
	go func() {
		result, err := a.transfer.GetByFromTo(ctx, &models.FilterTransfer{
			ProviderID: strconv.FormatInt(resultUser.SupId, 10),
			From:       input.From,
			To:         input.To,
		})
		if err != nil {
			errVar = err
		}
		resultTransfer = result
		wg.Done()
	}()

	wg.Wait()

	if errVar != nil {
		return nil, errVar
	}

	return &TransactionHistoryOutput{
		Deposit:  resultDeposit,
		Transfer: resultTransfer,
	}, nil
}
func (a *Usecase) Update(ctx context.Context, tx *gorm.DB, inputWhere *models.Member, inputUpdate map[string]interface{}) error {
	return a.uow.Member.Update(ctx, tx, inputWhere, inputUpdate)
}
func (a *Usecase) UpdatePassword(ctx context.Context, input *UpdatePasswordInput) error {
	result, _ := a.GetById(ctx, input.Id)
	if result == nil || !result.ComparePassword(input.Password) {
		return ErrIdOrPasswordInvalid
	}
	updateM := models.Member{
		Password: input.NewPassword,
	}
	updateM.HashPassword()

	if err := a.uow.Member.Update(ctx, nil, result, map[string]interface{}{"Password": updateM.Password, "Salt": updateM.Salt}); err != nil {
		return err
	}
	return nil
}

func (a *Usecase) Create(ctx context.Context, input *CreateInput) error {
	if err := utils.Validate.Struct(input); err != nil {
		return commons.ErrInvalidRequest(err)
	}
	var wg sync.WaitGroup
	var supIdGetResult *models.Member
	var emailGetResult *models.Member
	var err error
	if input.ReferralCode != nil && *input.ReferralCode != "" {
		wg.Add(1)
		go func() {
			_, errVar := a.Get(ctx, &models.Member{SupId: input.SupId})
			if errVar != nil {
				err = ErrReferralCode
			}
			wg.Done()
		}()
	}
	wg.Add(2)
	go func() {
		result, _ := a.Get(ctx, &models.Member{SupId: input.SupId})
		supIdGetResult = result
		wg.Done()
	}()
	go func() {
		result, _ := a.Get(ctx, &models.Member{Email: input.Email})
		emailGetResult = result
		wg.Done()
	}()

	wg.Wait()
	if err != nil {
		return err
	}
	if supIdGetResult != nil {
		return ErrIdExisted
	}
	if emailGetResult != nil {
		return ErrEmailExisted
	}

	member := &models.Member{
		SupId:        input.SupId,
		FullName:     input.FullName,
		Password:     input.Password,
		Email:        input.Email,
		Ip:           input.Ip,
		ReferralCode: input.ReferralCode,
		Status:       models.MemberStatusActive,
	}
	member.HashPassword()

	if err := a.uow.Member.Create(ctx, nil, member); err != nil {
		return err
	}

	return nil
}

func (u *Usecase) Deposit(ctx context.Context, id uint, input *DepositInput, mw middlewares.IMiddlewareManager) error {
	result, _ := u.GetById(ctx, id)
	if result == nil {
		return ErrNotFound
	}
	var resultReferMember *models.Member
	providerID, err := strconv.ParseInt(*result.ReferralCode, 10, 64)
	if err == nil {
		resultReferMember, _ = u.Get(ctx, &models.Member{SupId: providerID})
		if resultReferMember != nil {
			mw.SingleRequestAddUser(resultReferMember.ID)
			defer mw.SingleRequestReleaseUser(resultReferMember.ID)
		}
	}

	return u.uow.Do(ctx, func(tx *gorm.DB) error {
		memberDeposit, err := u.deposit.Post(ctx, tx, &deposit.CreateInput{
			SupId:         input.SupId,
			Num:           input.Num,
			ComissionType: input.ComissionType,
			Member:        result,
		})
		if err != nil {
			return err
		}

		if err := u.uow.Member.Update(ctx, tx, result, memberDeposit.AmountStrategy.UpdateMember(result, memberDeposit.Deposit)); err != nil {
			return err
		}
		if memberDeposit.AmountStrategy.IsDeposit() {
			if resultReferMember != nil {
				createComission := &comission.CreateInput{
					DepositId: memberDeposit.Deposit.ID,
					MemberId:  resultReferMember.ID,
					Amount:    memberDeposit.ComissionRefer,
				}

				if _, err := u.comission.Post(ctx, tx, createComission); err != nil {
					return err
				}
				updateMR := map[string]interface{}{
					"Commission": gorm.Expr(fmt.Sprintf("commission + %d", memberDeposit.ComissionRefer)),
				}

				if err := u.uow.Member.Update(ctx, tx, resultReferMember, updateMR); err != nil {
					return err
				}
			}
			if err := u.deposit.Deposit(ctx, memberDeposit); err != nil {
				return err
			}
		}

		return nil
	})
}

func (a *Usecase) UpdateMember(ctx context.Context, input *UpdateMember) error {
	result, err := a.GetById(ctx, input.Id)
	if err != nil {
		return err
	}

	if err := a.uow.Member.Update(ctx, nil, result, map[string]interface{}{"IsAgency": input.IsAgency, "Status": input.Status}); err != nil {
		return err
	}
	return nil
}
