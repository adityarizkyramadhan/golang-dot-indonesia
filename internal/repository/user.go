package repository

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/entity"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/infrastructure/cache"
	"gorm.io/gorm"
)

type User struct {
	db    *gorm.DB
	redis *cache.Redis
}

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id int) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int) error
}

func NewUser(db *gorm.DB, redis *cache.Redis) UserRepository {
	return &User{db, redis}
}

func (u *User) Create(ctx context.Context, user *entity.User) error {
	sqlStatement := `INSERT INTO users (username, password, name, role) VALUES (?, ?, ?, ?)`
	if err := u.db.WithContext(ctx).Exec(sqlStatement, user.Username, user.Password, user.Name, user.Role).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) FindByID(ctx context.Context, id int) (*entity.User, error) {
	// cek data di redis
	var user entity.User
	idStr := strconv.Itoa(id)
	key := "user:" + idStr
	if err := u.redis.Get(ctx, key, &user); err == nil {
		return &user, nil
	}
	if err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	go func() {
		expiredTime := time.Hour
		idStr := strconv.Itoa(id)
		key := "user:" + idStr
		if err := u.redis.Set(ctx, key, user, expiredTime); err != nil {
			log.Println(err)
		}
	}()
	return &user, nil
}

func (u *User) Update(ctx context.Context, user *entity.User) error {
	sqlStatement := `UPDATE users SET name = ?, role = ? WHERE id = ?`
	if err := u.db.WithContext(ctx).Exec(sqlStatement, user.Name, user.Role, user.ID).Error; err != nil {
		return err
	}

	go func() {
		idStr := strconv.Itoa(*user.ID)
		key := "user:" + idStr
		if err := u.redis.Del(ctx, key); err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func (u *User) Delete(ctx context.Context, id int) error {
	sqlStatement := `DELETE FROM users WHERE id = ?`
	if err := u.db.WithContext(ctx).Exec(sqlStatement, id).Error; err != nil {
		return err
	}

	go func() {
		idStr := strconv.Itoa(id)
		key := "user:" + idStr
		if err := u.redis.Del(ctx, key); err != nil {
			log.Println(err)
		}
	}()
	return nil
}
