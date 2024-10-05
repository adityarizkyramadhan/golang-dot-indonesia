package entity

type Goods struct {
	ID        *int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      *string  `gorm:"not null" json:"name"`
	Price     *float64 `gorm:"not null" json:"price"`
	Stock     *int     `gorm:"not null" json:"stock"`
	CreatedAt *string  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *string  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *string  `gorm:"null" json:"deleted_at"`
}

func (g *Goods) TableName() string {
	return "goods"
}
