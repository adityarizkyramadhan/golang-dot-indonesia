package entity

import (
	"time"

	"github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/generator"
)

type InvoicePurchase struct {
	ID            *int       `gorm:"primaryKey;autoIncrement" json:"id"`
	InvoiceNumber *string    `gorm:"unique;not null" json:"invoice_number"`
	OrderDate     *string    `gorm:"not null" json:"order_date"`
	IsPaid        *bool      `gorm:"not null" json:"is_paid"`
	Total         *float64   `gorm:"not null" json:"total"`
	CreatedAt     *string    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *string    `gorm:"autoUpdateTime" json:"updated_at"`
	Purchases     []Purchase `gorm:"foreignKey:InvoicePurchaseID" json:"purchases"`
}

func (i *InvoicePurchase) TableName() string {
	return "invoice_purchases"
}

func (i *InvoicePurchase) BeforeCreate() error {
	if i.IsPaid == nil {
		i.IsPaid = new(bool)
		*i.IsPaid = false
	}
	if i.InvoiceNumber == nil {
		inv := generator.InvoiceGenerator()
		i.InvoiceNumber = inv
	}
	if i.OrderDate == nil {
		now := time.Now().Format("2006-01-02")
		i.OrderDate = &now
	}
	return nil
}

type Purchase struct {
	ID                *int     `gorm:"primaryKey;autoIncrement" json:"id"`
	InvoicePurchaseID *int     `gorm:"not null" json:"invoice_purchase_id"`
	GoodsID           *int     `gorm:"not null" json:"goods_id"`
	Amount            *int     `gorm:"not null" json:"amount"`
	PricePerUnit      *float64 `gorm:"not null" json:"price_per_unit"`
	TotalPrice        *float64 `gorm:"not null" json:"total_price"`
	CreatedAt         *string  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         *string  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Purchase) TableName() string {
	return "purchases"
}
