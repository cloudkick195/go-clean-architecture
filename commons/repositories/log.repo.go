package repositories

import (
	"context"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"

	"gorm.io/gorm"
)

// Interface logRepository
type ILogRepo interface {
	IBase
	Create(ctx context.Context, m *models.Log) error
	Get(ctx context.Context, m *models.Log) (*models.Log, error)
	GetById(ctx context.Context, id uint) (*models.Log, error)
	Updates(ctx context.Context, model *models.Log, Id uint, updates map[string]interface{}) (err error)
}

// Struct logRepository
type logRepo struct {
	*base
}

func NewLogRepo(dbConnect *gorm.DB) ILogRepo {
	newModel := models.Log{}
	return &logRepo{base: &base{db: dbConnect, table: newModel.TableName()}}
}

func (r *logRepo) Create(ctx context.Context, m *models.Log) error {
	err := r.db.WithContext(ctx).Create(&m).Error

	if err != nil {
		return commons.ErrDB(err)
	}
	return nil
}

func (r *logRepo) Get(ctx context.Context, input *models.Log) (*models.Log, error) {
	err := r.db.WithContext(ctx).Where(input).First(&input).Error

	if err != nil {

		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *logRepo) GetById(ctx context.Context, id uint) (*models.Log, error) {
	var input *models.Log
	err := r.db.WithContext(ctx).First(&input, id).Error

	if err != nil {
		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *logRepo) Updates(ctx context.Context, model *models.Log, Id uint, updates map[string]interface{}) (err error) {
	return r.db.WithContext(ctx).Model(model).Where("id = ?", Id).Updates(updates).Error
}
