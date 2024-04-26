package repositories

import (
	"context"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"

	"gorm.io/gorm"
)

// Interface depositLogRepository
type IDepositLogRepo interface {
	IBase
	Create(ctx context.Context, m *models.DepositLog) error
	Get(ctx context.Context, m *models.DepositLog) (*models.DepositLog, error)
	GetById(ctx context.Context, id uint) (*models.DepositLog, error)
	Updates(ctx context.Context, model *models.DepositLog, Id uint, updates map[string]interface{}) (err error)
	Save(ctx context.Context, m *models.DepositLog) error
}

// Struct depositLogRepository
type depositLogRepo struct {
	*base
}

func NewDepositLogRepo(dbConnect *gorm.DB) IDepositLogRepo {
	newModel := models.DepositLog{}
	return &depositLogRepo{base: &base{db: dbConnect, table: newModel.TableName()}}
}

func (r *depositLogRepo) Create(ctx context.Context, m *models.DepositLog) error {
	err := r.db.WithContext(ctx).Create(&m).Error

	if err != nil {
		return commons.ErrDB(err)
	}
	return nil
}

func (r *depositLogRepo) Get(ctx context.Context, input *models.DepositLog) (*models.DepositLog, error) {
	err := r.db.WithContext(ctx).Where(input).First(&input).Error

	if err != nil {

		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *depositLogRepo) GetById(ctx context.Context, id uint) (*models.DepositLog, error) {
	var input *models.DepositLog
	err := r.db.WithContext(ctx).First(&input, id).Error

	if err != nil {

		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *depositLogRepo) Updates(ctx context.Context, model *models.DepositLog, Id uint, updates map[string]interface{}) error {
	err := r.db.WithContext(ctx).Model(model).Where("id = ?", Id).Updates(updates).Error
	if err != nil {
		return commons.ErrDB(err)
	}
	return nil
}

func (r *depositLogRepo) Save(ctx context.Context, m *models.DepositLog) error {
	err := r.db.WithContext(ctx).Save(m).Error

	if err != nil {
		return commons.ErrDB(err)
	}

	return nil
}
