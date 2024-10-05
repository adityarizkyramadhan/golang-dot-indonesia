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

func NewPurchase(db *gorm.DB, redis *cache.Redis, repoGood *Good) PurchaseRepository {
	return &Purchase{db: db, redis: redis, repoGood: repoGood}
}

// Modify the PurchaseRepository interface
type PurchaseRepository interface {
	Create(ctx context.Context, invoicePurchase *entity.InvoicePurchase, purchases []entity.Purchase) error
	Get(ctx context.Context, tx *gorm.DB, id int, userID *int) (*entity.InvoicePurchase, error)          // Added userID
	Update(ctx context.Context, tx *gorm.DB, invoicePurchase *entity.InvoicePurchase, userID *int) error // Added userID
	Delete(ctx context.Context, tx *gorm.DB, id int, userID *int) error                                  // Added userID
	GetAll(ctx context.Context, tx *gorm.DB, userID *int) ([]entity.InvoicePurchase, error)              // Added userID
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

// Modify the Purchase struct to include the necessary methods
func (r *Purchase) Get(ctx context.Context, tx *gorm.DB, id int, userID *int) (*entity.InvoicePurchase, error) {
	var invoicePurchase entity.InvoicePurchase
	if tx == nil {
		tx = r.db
	}
	// Add user ID check if required
	query := tx.WithContext(ctx).Preload("Purchases").First(&invoicePurchase, id)
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	err := query.Error
	if err != nil {
		return nil, err
	}
	return &invoicePurchase, nil
}

func (r *Purchase) Update(ctx context.Context, tx *gorm.DB, invoicePurchase *entity.InvoicePurchase, userID *int) error {
	if tx == nil {
		tx = r.db
	}
	// Optionally check user ID
	if userID != nil {
		// Add logic here to ensure the invoice belongs to the user
	}
	return tx.WithContext(ctx).Save(invoicePurchase).Error
}

func (r *Purchase) Delete(ctx context.Context, tx *gorm.DB, id int, userID *int) error {
	if tx == nil {
		tx = r.db
	}
	var invoicePurchase entity.InvoicePurchase
	err := tx.WithContext(ctx).First(&invoicePurchase, id).Error
	if err != nil {
		return err
	}
	// Optionally check if userID matches
	return tx.WithContext(ctx).Delete(&invoicePurchase).Error
}

func (r *Purchase) GetAll(ctx context.Context, tx *gorm.DB, userID *int) ([]entity.InvoicePurchase, error) {
	var invoicePurchases []entity.InvoicePurchase
	if tx == nil {
		tx = r.db
	}
	// Optionally add filtering by user ID
	query := tx.WithContext(ctx).Preload("Purchases")
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	err := query.Find(&invoicePurchases).Error
	if err != nil {
		return nil, err
	}
	return invoicePurchases, nil
}
