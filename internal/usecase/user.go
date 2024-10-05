package usecase

import (
	"context"

	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/dto"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/entity"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/repository"
	custom_error "github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/errors"
	utils_jwt "github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	userRepo repository.UserRepository
}

type UserUsecase interface {
	Create(ctx context.Context, user *dto.UserRegister) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id int) (*entity.User, error)
	Login(ctx context.Context, user *dto.UserLogin) (string, error)
	Update(ctx context.Context, user *dto.UserUpdate) error
	Delete(ctx context.Context, id int) error
}

func NewUser(userRepo repository.UserRepository) UserUsecase {
	return &User{userRepo}
}

func (u *User) Create(ctx context.Context, user *dto.UserRegister) error {
	role := "user"
	password := user.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return custom_error.NewError(custom_error.ErrInternalServer, err.Error())
	}
	hashedPasswordString := string(hashedPassword)
	userEntity := &entity.User{
		Username: user.Username,
		Password: &hashedPasswordString,
		Name:     user.Name,
		Role:     &role,
	}
	return u.userRepo.Create(ctx, userEntity)
}

func (u *User) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	return u.userRepo.FindByUsername(ctx, username)
}

func (u *User) FindByID(ctx context.Context, id int) (*entity.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *User) Login(ctx context.Context, user *dto.UserLogin) (string, error) {
	userEntity, err := u.FindByUsername(ctx, *user.Username)
	if err != nil {
		return "", custom_error.NewError(custom_error.ErrNotFound, "user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*userEntity.Password), []byte(*user.Password))
	if err != nil {
		return "", custom_error.NewError(custom_error.ErrUnauthorized, "invalid password")
	}

	token, err := utils_jwt.GenerateToken(*userEntity.ID)
	if err != nil {
		return "", custom_error.NewError(custom_error.ErrInternalServer, err.Error())
	}

	return token, nil
}

func (u *User) Update(ctx context.Context, user *dto.UserUpdate) error {
	userEntity, err := u.FindByID(ctx, *user.ID)
	if err != nil {
		return custom_error.NewError(custom_error.ErrNotFound, "user not found")
	}

	userEntity.Name = user.Name
	return u.userRepo.Update(ctx, userEntity)
}

func (u *User) Delete(ctx context.Context, id int) error {
	userEntity, err := u.FindByID(ctx, id)
	if err != nil {
		return custom_error.NewError(custom_error.ErrNotFound, "user not found")
	}

	return u.userRepo.Delete(ctx, *userEntity.ID)
}
