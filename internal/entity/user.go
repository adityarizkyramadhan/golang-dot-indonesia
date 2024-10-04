package entity

type User struct {
	ID       *int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username *string `gorm:"unique;not null" json:"username"`
	Password *string `gorm:"not null" json:"-"`
	Name     *string `gorm:"not null;size:255" json:"name"`
	Role     *string `gorm:"not null;size:10;default:'user'" json:"role"`
}

func (u *User) TableName() string {
	return "users"
}
