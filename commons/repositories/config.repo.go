package repositories

import (
	"context"
	"errors"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"

	"gorm.io/gorm"
)

// Interface configRepository
type IConfigRepo interface {
	Get(ctx context.Context, m *models.Config) (*models.Config, error)
	List(ctx context.Context) ([]models.Config, error)
	IBase
}

// Struct configRepository
type configRepo struct {
	*base
}

func NewConfigRepo(dbConnect *gorm.DB) IConfigRepo {
	newModel := models.Config{}
	return &configRepo{base: &base{db: dbConnect, table: newModel.TableName()}}
}

func (r *configRepo) Create(ctx context.Context, m *models.Config) error {
	err := r.db.WithContext(ctx).Create(&m).Error

	if err != nil {
		return commons.ErrDB(err)
	}
	return nil
}

func (r *configRepo) List(ctx context.Context) (list []models.Config, err error) {
	err = r.db.WithContext(ctx).Find(&list).Error
	return list, err
}

func (r *configRepo) Get(ctx context.Context, input *models.Config) (*models.Config, error) {
	err := r.db.WithContext(ctx).Where(input).First(&input).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *configRepo) GetById(ctx context.Context, id uint) (*models.Config, error) {
	var input *models.Config
	err := r.db.WithContext(ctx).First(&input, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *configRepo) Updates(ctx context.Context, model *models.Config, Id uint, updates map[string]interface{}) (err error) {
	return r.db.WithContext(ctx).Model(model).Where("id = ?", Id).Updates(updates).Error
}
