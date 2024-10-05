package entity

import "time"

type Goods struct {
	ID        *int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      *string  `gorm:"not null" json:"name"`
	Price     *float64 `gorm:"not null" json:"price"`
	Stock     *int     `gorm:"not null" json:"stock"`
	CreatedAt *string  `json:"created_at"`
	UpdatedAt *string  `json:"updated_at"`
}

func (g *Goods) TableName() string {
	return "goods"
}

func (g *Goods) BeforeCreate() error {
	now := time.Now().Format(time.RFC3339)
	g.CreatedAt = &now
	g.UpdatedAt = &now
	return nil
}

func (g *Goods) BeforeUpdate() error {
	now := time.Now().Format(time.RFC3339)
	g.UpdatedAt = &now
	return nil
}
