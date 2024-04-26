package uows

import (
	"context"
	"go_clean_architecture/commons/repositories"

	"gorm.io/gorm"
)

type IRepoUOW interface {
	Begin()
	Rollback()
	Commit()
}

type RepoUOW struct {
	db          *gorm.DB
	Config      repositories.IConfigRepo
	Deposit     repositories.IDepositRepo
	DepositLog  repositories.IDepositLogRepo
	Log         repositories.ILogRepo
	Member      repositories.IMemberRepo
	Transfer    repositories.ITransferRepo
	TransferLog repositories.ITransferLogRepo
	Comission   repositories.IComissionRepo
}

func NewRepoUnitOfWork(db *gorm.DB) *RepoUOW {
	return &RepoUOW{
		db:          db,
		Config:      repositories.NewConfigRepo(db),
		Deposit:     repositories.NewDepositRepo(db),
		DepositLog:  repositories.NewDepositLogRepo(db),
		Log:         repositories.NewLogRepo(db),
		Member:      repositories.NewMemberRepo(db),
		Transfer:    repositories.NewTransferRepo(db),
		TransferLog: repositories.NewTransferLogRepo(db),
		Comission:   repositories.NewComissionRepo(db),
	}
}
func (uow *RepoUOW) Do(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := uow.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
