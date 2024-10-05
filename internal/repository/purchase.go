package repository

import (
	"context"

	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/entity"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/infrastructure/cache"
	custom_error "github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/errors"
	"gorm.io/gorm"
)

type Purchase struct {
	db       *gorm.DB
	redis    *cache.Redis
	repoGood *Good
}

type PurchaseRepository interface {
	Create(ctx context.Context, invoicePurchase *entity.InvoicePurchase, purchases []entity.Purchase) error
	Get(ctx context.Context, tx *gorm.DB, id int) (*entity.InvoicePurchase, error)
	Update(ctx context.Context, tx *gorm.DB, invoicePurchase *entity.InvoicePurchase) error
	Delete(ctx context.Context, tx *gorm.DB, id int) error
	GetAll(ctx context.Context, tx *gorm.DB) ([]entity.InvoicePurchase, error)
}

func NewPurchase(db *gorm.DB, redis *cache.Redis, repoGood *Good) PurchaseRepository {
	return &Purchase{db: db, redis: redis, repoGood: repoGood}
}

func (r *Purchase) Create(ctx context.Context, invoicePurchase *entity.InvoicePurchase, purchases []entity.Purchase) error {
	tx := r.db.WithContext(ctx).Begin()
	invoicePurchase.BeforeCreate()
	total := 0.0
	for _, purchase := range purchases {
		total += *purchase.TotalPrice
	}
	invoicePurchase.Total = &total
	if err := tx.Create(invoicePurchase).Error; err != nil {
		tx.Rollback()
		return custom_error.NewError(custom_error.ErrInternalServer, err.Error())
	}

	for i := range purchases {
		purchases[i].InvoicePurchaseID = invoicePurchase.ID
		purchases[i].BeforeCreate()
	}

	if err := tx.Create(&purchases).Error; err != nil {
		tx.Rollback()
		return custom_error.NewError(custom_error.ErrInternalServer, err.Error())
	}

	for _, purchase := range purchases {
		if err := r.repoGood.DecrementStock(ctx, tx, *purchase.GoodsID, *purchase.Amount); err != nil {
			tx.Rollback()
			return custom_error.NewError(custom_error.ErrInternalServer, err.Error())
		}
	}

	tx.Commit()
	return nil
}

func (r *Purchase) Get(ctx context.Context, tx *gorm.DB, id int) (*entity.InvoicePurchase, error) {
	var invoicePurchase entity.InvoicePurchase
	if tx == nil {
		tx = r.db
	}
	err := tx.WithContext(ctx).Preload("Purchases").First(&invoicePurchase, id).Error
	if err != nil {
		return nil, err
	}
	return &invoicePurchase, nil
}

func (r *Purchase) Update(ctx context.Context, tx *gorm.DB, invoicePurchase *entity.InvoicePurchase) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).Save(invoicePurchase).Error
}

func (r *Purchase) Delete(ctx context.Context, tx *gorm.DB, id int) error {
	if tx == nil {
		tx = r.db
	}
	var invoicePurchase entity.InvoicePurchase
	err := tx.WithContext(ctx).First(&invoicePurchase, id).Error
	if err != nil {
		return err
	}
	return tx.WithContext(ctx).Delete(&invoicePurchase).Error
}

func (r *Purchase) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.InvoicePurchase, error) {
	var invoicePurchases []entity.InvoicePurchase
	if tx == nil {
		tx = r.db
	}
	err := tx.WithContext(ctx).Preload("Purchases").Find(&invoicePurchases).Error
	if err != nil {
		return nil, err
	}
	return invoicePurchases, nil
}
