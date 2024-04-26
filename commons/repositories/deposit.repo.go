package repositories

import (
	"context"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"

	"gorm.io/gorm"
)

// Interface depositRepository
type IDepositRepo interface {
	IBase
	Create(ctx context.Context, tx *gorm.DB, m *models.Deposit) error
	Get(ctx context.Context, m *models.Deposit) (*models.Deposit, error)
	GetById(ctx context.Context, id uint) (*models.Deposit, error)
	GetByFromTo(ctx context.Context, input *models.FilterDeposit) (result []models.Deposit, err error)
	List(ctx context.Context, paging *commons.Pagination, queries *models.FilterDeposit) (results []models.Deposit, err error)
	Update(ctx context.Context, tx *gorm.DB, record *models.Deposit, updates map[string]interface{}) error
}

// Struct depositRepository
type depositRepo struct {
	*base
}

func NewDepositRepo(dbConnect *gorm.DB) IDepositRepo {
	newModel := models.Deposit{}
	return &depositRepo{base: &base{db: dbConnect, table: newModel.TableName()}}
}

func (r *depositRepo) Create(ctx context.Context, tx *gorm.DB, m *models.Deposit) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.WithContext(ctx).Create(&m).Error; err != nil {
		return commons.ErrDB(err)
	}
	return nil
}

func (r *depositRepo) Get(ctx context.Context, input *models.Deposit) (*models.Deposit, error) {
	err := r.db.WithContext(ctx).Where(input).First(&input).Error

	if err != nil {

		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *depositRepo) Update(ctx context.Context, tx *gorm.DB, record *models.Deposit, updates map[string]interface{}) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	if err := db.WithContext(ctx).Model(record).Updates(updates).Error; err != nil {
		return commons.ErrDB(err)
	}
	return nil
}

func (r *depositRepo) List(ctx context.Context, paging *commons.Pagination, queries *models.FilterDeposit) (results []models.Deposit, err error) {
	query := r.db.WithContext(ctx).Preload("Member")

	if queries.ComissionType != nil {
		query = query.Where("comission_type = ?", queries.ComissionType)
	}
	if queries.MemberId != nil {
		query = query.Where("provider_id = ?", queries.MemberId)
	}

	query.Find(&results).Count(&paging.Total)
	err = query.Limit(paging.Limit).Offset(paging.Offset).Order(paging.Sort).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *depositRepo) GetById(ctx context.Context, id uint) (*models.Deposit, error) {
	var input *models.Deposit
	err := r.db.WithContext(ctx).First(&input, id).Error

	if err != nil {
		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *depositRepo) GetByFromTo(ctx context.Context, input *models.FilterDeposit) (result []models.Deposit, err error) {
	from := input.From
	to := input.To
	input.From = ""
	input.To = ""
	err = r.db.WithContext(ctx).Order("created_at desc").
		Where(input).
		Where("Created_at BETWEEN ? AND ?", from, to).Find(&result).Error
	if err != nil {
		return nil, commons.ErrDB(err)
	}
	return result, err
}
