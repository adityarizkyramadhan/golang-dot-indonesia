package entity

import "time"

type User struct {
	ID        *int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  *string `gorm:"unique;not null" json:"username"`
	Password  *string `gorm:"not null" json:"-"`
	Name      *string `gorm:"not null;size:255" json:"name"`
	Role      *string `gorm:"not null;size:10;default:'user'" json:"role"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate() error {
	now := time.Now().Format(time.RFC3339)
	u.CreatedAt = &now
	u.UpdatedAt = &now
	return nil
}

func (u *User) BeforeUpdate() error {
	now := time.Now().Format(time.RFC3339)
	u.UpdatedAt = &now
	return nil
}
