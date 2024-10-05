package repository

import (
	"context"
	"sync"

	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/entity"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/infrastructure/cache"
	custom_error "github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/errors"
	"gorm.io/gorm"
)

type Good struct {
	db    *gorm.DB
	mutex sync.Mutex
	redis *cache.Redis
}

func NewGoods(db *gorm.DB, redis *cache.Redis) *Good {
	good := &Good{
		db:    db,
		mutex: sync.Mutex{},
		redis: redis,
	}
	return good
}

func (r *Good) Create(ctx context.Context, tx *gorm.DB, good *entity.Goods) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).Create(good).Error
}

func (r *Good) Get(ctx context.Context, tx *gorm.DB, id int) (*entity.Goods, error) {
	var good entity.Goods
	if tx == nil {
		tx = r.db
	}
	err := tx.WithContext(ctx).First(&good, id).Error
	if err != nil {
		return nil, err
	}
	return &good, nil
}

func (r *Good) Update(ctx context.Context, tx *gorm.DB, good *entity.Goods) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).Save(good).Error
}

func (r *Good) Delete(ctx context.Context, tx *gorm.DB, id int) error {
	if tx == nil {
		tx = r.db
	}
	var good entity.Goods
	err := tx.WithContext(ctx).First(&good, id).Error
	if err != nil {
		return err
	}
	return tx.WithContext(ctx).Delete(&good).Error
}

func (r *Good) DecrementStock(ctx context.Context, tx *gorm.DB, id int, amount int) error {
	if tx == nil {
		tx = r.db
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var good entity.Goods
	err := tx.WithContext(ctx).First(&good, id).Error
	if err != nil {
		return err
	}

	if good.Stock == nil || *good.Stock < amount {
		return custom_error.NewError(custom_error.ErrBadRequest, "insufficient stock")
	}

	newStock := *good.Stock - amount
	good.Stock = &newStock

	return tx.WithContext(ctx).Save(&good).Error
}

func (r *Good) AddStock(ctx context.Context, tx *gorm.DB, id int, amount int) error {
	if tx == nil {
		tx = r.db
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var good entity.Goods
	err := tx.WithContext(ctx).First(&good, id).Error
	if err != nil {
		return err
	}

	newStock := *good.Stock + amount
	good.Stock = &newStock

	err = tx.WithContext(ctx).Save(&good).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Good) List(ctx context.Context, name *string, page int, pageSize int) ([]entity.Goods, int64, error) {
	var goods []entity.Goods
	var total int64

	query := r.db.Model(&entity.Goods{})

	if name != nil {
		query = query.Where("name LIKE ?", "%"+*name+"%")
	}

	if err := query.WithContext(ctx).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&goods).Error
	if err != nil {
		return nil, 0, err
	}

	return goods, total, nil
}
