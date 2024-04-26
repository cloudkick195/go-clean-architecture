package repositories

import (
	"context"
	"errors"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"

	"gorm.io/gorm"
)

// Interface transferRepository
type ITransferRepo interface {
	IBase

	Get(ctx context.Context, m *models.Transfer) (*models.Transfer, error)
	GetById(ctx context.Context, id uint) (*models.Transfer, error)
	GetByFromTo(ctx context.Context, input *models.FilterTransfer) (result []models.Transfer, err error)
	FirstOrCreate(ctx context.Context, model *models.Transfer) error
	Save(ctx context.Context, m *models.Transfer) error
	List(ctx context.Context, paging *commons.Pagination, queries *models.FilterTransfer) (results []models.Transfer, err error)
}

// Struct transferRepository
type transferRepo struct {
	*base
}

func NewTransferRepo(dbConnect *gorm.DB) ITransferRepo {
	newModel := models.Transfer{}
	return &transferRepo{base: &base{db: dbConnect, table: newModel.TableName()}}
}

func (r *transferRepo) List(ctx context.Context, paging *commons.Pagination, queries *models.FilterTransfer) (results []models.Transfer, err error) {
	query := r.db.WithContext(ctx)

	if queries.ProviderID != "" {
		query = query.Where("provider_id = ?", queries.ProviderID)
	}

	query.Find(&results).Count(&paging.Total)
	err = query.Limit(paging.Limit).Offset(paging.Offset).Order(paging.Sort).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *transferRepo) GetByFromTo(ctx context.Context, input *models.FilterTransfer) (result []models.Transfer, err error) {
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

func (r *transferRepo) Create(ctx context.Context, m *models.Transfer) error {
	result := r.db.WithContext(ctx).Create(&m)

	if result.Error != nil {
		return commons.ErrDB(result.Error)
	}
	return nil
}

func (r *transferRepo) FirstOrCreate(ctx context.Context, model *models.Transfer) error {
	return r.db.WithContext(ctx).Where(models.Transfer{PrivateID: model.PrivateID}).FirstOrCreate(model).Error
}

func (r *transferRepo) Get(ctx context.Context, input *models.Transfer) (*models.Transfer, error) {
	err := r.db.WithContext(ctx).Where(input).First(&input).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *transferRepo) GetById(ctx context.Context, id uint) (*models.Transfer, error) {
	var input *models.Transfer
	err := r.db.WithContext(ctx).First(&input, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *transferRepo) Save(ctx context.Context, m *models.Transfer) error {
	err := r.db.WithContext(ctx).Save(m).Error

	if err != nil {
		return commons.ErrDB(err)
	}

	return nil
}
