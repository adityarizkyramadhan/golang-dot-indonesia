package handler

import (
	"net/http"

	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/dto"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/middleware"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/usecase"
	custom_error "github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/errors"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/response"
	"github.com/gin-gonic/gin"
)

type User struct {
	userUsece usecase.UserUsecase
}

func NewUser(userUsece usecase.UserUsecase) *User {
	return &User{userUsece}
}

func (u *User) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/register", u.Register)
	r.POST("/login", u.Login)
	r.GET("/profile", middleware.JWTMiddleware(), u.Profile)
}

func (u *User) Register(ctx *gin.Context) {
	var user dto.UserRegister
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(custom_error.NewError(custom_error.ErrBadRequest, err.Error()))
		ctx.Next()
		return
	}

	err := u.userUsece.Create(ctx, &user)
	if err != nil {
		ctx.Error(err)
		ctx.Next()
		return
	}

	response := response.Success("user created")
	ctx.JSON(http.StatusCreated, response)
}

func (u *User) Login(ctx *gin.Context) {
	var user dto.UserLogin
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(custom_error.NewError(custom_error.ErrBadRequest, err.Error()))
		ctx.Next()
		return
	}

	token, err := u.userUsece.Login(ctx, &user)
	if err != nil {
		ctx.Error(err)
		ctx.Next()
		return
	}

	response := response.Success(gin.H{"token": token})
	ctx.JSON(http.StatusOK, response)
}

func (u *User) Profile(ctx *gin.Context) {
	id := ctx.MustGet("id").(int)
	user, err := u.userUsece.FindByID(ctx, id)
	if err != nil {
		ctx.Error(err)
		ctx.Next()
		return
	}
	response := response.Success(user)
	ctx.JSON(http.StatusOK, response)
}
