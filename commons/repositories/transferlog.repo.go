package repositories

import (
	"context"
	"errors"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"

	"gorm.io/gorm"
)

// Interface transferLogRepository
type ITransferLogRepo interface {
	IBase
	Create(ctx context.Context, m *models.TransferLog) error
	Get(ctx context.Context, m *models.TransferLog) (*models.TransferLog, error)
	GetById(ctx context.Context, id uint) (*models.TransferLog, error)
	Updates(ctx context.Context, model *models.TransferLog, Id uint, updates map[string]interface{}) (err error)
	Save(ctx context.Context, m *models.TransferLog) error
}

// Struct transferLogRepository
type transferLogRepo struct {
	*base
}

func NewTransferLogRepo(dbConnect *gorm.DB) ITransferLogRepo {
	newModel := models.TransferLog{}
	return &transferLogRepo{base: &base{db: dbConnect, table: newModel.TableName()}}
}
func (r *transferLogRepo) Create(ctx context.Context, m *models.TransferLog) error {
	err := r.db.WithContext(ctx).Create(&m).Error

	if err != nil {
		return commons.ErrDB(err)
	}
	return nil
}

func (r *transferLogRepo) Get(ctx context.Context, input *models.TransferLog) (*models.TransferLog, error) {
	err := r.db.WithContext(ctx).Where(input).First(&input).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *transferLogRepo) GetById(ctx context.Context, id uint) (*models.TransferLog, error) {
	var input *models.TransferLog
	err := r.db.WithContext(ctx).First(&input, id).Error

	if err != nil {
		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *transferLogRepo) Updates(ctx context.Context, model *models.TransferLog, Id uint, updates map[string]interface{}) error {
	err := r.db.WithContext(ctx).Model(model).Where("id = ?", Id).Updates(updates).Error
	if err != nil {
		return commons.ErrDB(err)
	}
	return nil
}

func (r *transferLogRepo) Save(ctx context.Context, m *models.TransferLog) error {
	err := r.db.WithContext(ctx).Save(m).Error

	if err != nil {
		return commons.ErrDB(err)
	}

	return nil
}
