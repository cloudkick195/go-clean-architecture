package repositories

import (
	"context"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"

	"gorm.io/gorm"
)

// Interface ComissionRepository
type IComissionRepo interface {
	IBase
	Create(ctx context.Context, tx *gorm.DB, m *models.Comission) error
	Get(ctx context.Context, m *models.Comission) (*models.Comission, error)
	GetById(ctx context.Context, id uint) (*models.Comission, error)
}

// Struct ComissionRepository
type ComissionRepo struct {
	*base
}

func NewComissionRepo(dbConnect *gorm.DB) IComissionRepo {
	newModel := models.Comission{}
	return &ComissionRepo{base: &base{db: dbConnect, table: newModel.TableName()}}
}

func (r *ComissionRepo) Create(ctx context.Context, tx *gorm.DB, m *models.Comission) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.WithContext(ctx).Create(&m).Error; err != nil {
		return commons.ErrDB(err)
	}
	return nil
}

func (r *ComissionRepo) Get(ctx context.Context, input *models.Comission) (*models.Comission, error) {
	err := r.db.WithContext(ctx).Where(input).First(&input).Error

	if err != nil {

		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *ComissionRepo) GetById(ctx context.Context, id uint) (*models.Comission, error) {
	var input *models.Comission
	err := r.db.WithContext(ctx).First(&input, id).Error

	if err != nil {

		return nil, commons.ErrDB(err)
	}

	return input, nil
}
