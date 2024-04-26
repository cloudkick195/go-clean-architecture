package repositories

import (
	"context"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"

	"gorm.io/gorm"
)

type IMemberRepo interface {
	IBase
	Create(ctx context.Context, tx *gorm.DB, member *models.Member) error
	Update(ctx context.Context, tx *gorm.DB, member *models.Member, updates map[string]interface{}) error
	List(ctx context.Context, paging *commons.Pagination, queries *models.FilterMember) (results []models.Member, err error)
	Get(ctx context.Context, m *models.Member) (*models.Member, error)
	GetById(ctx context.Context, id uint) (*models.Member, error)
}

type memberRepo struct {
	*base
}

func NewMemberRepo(dbConnect *gorm.DB) IMemberRepo {
	newModel := models.Member{}
	return &memberRepo{base: &base{db: dbConnect, table: newModel.TableName()}}
}

func (r *memberRepo) Create(ctx context.Context, tx *gorm.DB, input *models.Member) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.WithContext(ctx).Create(&input).Error; err != nil {
		return err
	}
	return nil
}

func (r *memberRepo) Get(ctx context.Context, input *models.Member) (*models.Member, error) {
	err := r.db.WithContext(ctx).Where(input).First(&input).Error

	if err != nil {
		return nil, commons.ErrDB(err)
	}

	return input, nil
}

func (r *memberRepo) List(ctx context.Context, paging *commons.Pagination, queries *models.FilterMember) (results []models.Member, err error) {
	query := r.db.WithContext(ctx).Where("role is null")

	if paging.Query != "" {
		querySearch := "%" + paging.Query + "%"
		query = query.Where("sup_id = ? or full_name like ? or email like ?", paging.Query, querySearch, querySearch)
	}

	query.Find(&results).Count(&paging.Total)
	err = query.Limit(paging.Limit).Offset(paging.Offset).Order(paging.Sort).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *memberRepo) GetById(ctx context.Context, input uint) (*models.Member, error) {
	var member models.Member
	err := r.db.WithContext(ctx).First(&member, input).Error

	if err != nil {
		return nil, err
	}

	return &member, nil
}

// func (r *memberRepo) Update(ctx context.Context, member *models.Member, updates models.Member) error {
// 	if err := r.db.WithContext(ctx).Model(&member).Updates(updates).Error; err != nil {
// 		return commons.ErrDB(err)
// 	}
// 	return nil
// }

func (r *memberRepo) Update(ctx context.Context, tx *gorm.DB, member *models.Member, updates map[string]interface{}) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	if err := db.WithContext(ctx).Model(member).Updates(updates).Error; err != nil {
		return commons.ErrDB(err)
	}
	return nil
}
