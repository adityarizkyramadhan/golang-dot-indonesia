package handler

import (
	"net/http"
	"strconv"

	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/dto"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/entity"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/middleware"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/usecase"
	custom_error "github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/errors"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/response"
	"github.com/gin-gonic/gin"
)

type GoodsHandler struct {
	goodsUsecase usecase.GoodsUsecase
}

// NewGoodsHandler initializes a new GoodsHandler
func NewGoods(goodsUsecase usecase.GoodsUsecase) *GoodsHandler {
	return &GoodsHandler{goodsUsecase: goodsUsecase}
}

func (h *GoodsHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/", middleware.JWTMiddleware(), h.Create)
	r.GET("/:id", middleware.JWTMiddleware(), h.Get)
	r.PUT("/:id", middleware.JWTMiddleware(), h.Update)
	r.DELETE("/:id", middleware.JWTMiddleware(), h.Delete)
	r.PATCH("/:id/stock", middleware.JWTMiddleware(), h.AddStock)
	r.GET("/", middleware.JWTMiddleware(), h.List)
}

// Create a new good
func (h *GoodsHandler) Create(ctx *gin.Context) {
	var good entity.Goods
	if err := ctx.ShouldBindJSON(&good); err != nil {
		err := custom_error.NewError(custom_error.ErrBadRequest, err.Error())
		ctx.Error(err)
		ctx.Next()
		return
	}

	if err := h.goodsUsecase.Create(ctx.Request.Context(), &good); err != nil {
		ctx.Error(err)
		ctx.Next()
		return
	}

	response := response.Success(good)
	ctx.JSON(http.StatusCreated, response)
}

// Get retrieves a good by ID
func (h *GoodsHandler) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		err := custom_error.NewError(custom_error.ErrBadRequest, "invalid ID")
		ctx.Error(err)
		ctx.Next()
		return
	}
	good, err := h.goodsUsecase.Get(ctx.Request.Context(), idInt)
	if err != nil {
		ctx.Error(err)
		ctx.Next()
		return
	}

	response := response.Success(good)
	ctx.JSON(http.StatusOK, response)
}

// Update modifies an existing good
func (h *GoodsHandler) Update(c *gin.Context) {
	var good entity.Goods
	if err := c.ShouldBindJSON(&good); err != nil {
		err := custom_error.NewError(custom_error.ErrBadRequest, err.Error())
		c.Error(err)
		c.Next()
		return
	}

	if err := h.goodsUsecase.Update(c.Request.Context(), &good); err != nil {
		c.Error(err)
		c.Next()
		return
	}

	response := response.Success("good updated")
	c.JSON(http.StatusOK, response)
}

// Delete removes a good by ID
func (h *GoodsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		err := custom_error.NewError(custom_error.ErrBadRequest, "invalid ID")
		c.Error(err)
		c.Next()
		return
	}

	if err := h.goodsUsecase.Delete(c.Request.Context(), idInt); err != nil {
		c.Error(err)
		c.Next()
		return
	}

	response := response.Success("good deleted")
	c.JSON(http.StatusOK, response)
}

// AddStock increases the stock of a good
func (h *GoodsHandler) AddStock(c *gin.Context) {
	var stockUpdate struct {
		Amount int `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&stockUpdate); err != nil {
		err := custom_error.NewError(custom_error.ErrBadRequest, err.Error())
		c.Error(err)
		c.Next()
		return
	}

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		err := custom_error.NewError(custom_error.ErrBadRequest, "invalid ID")
		c.Error(err)
		c.Next()
		return
	}

	if err := h.goodsUsecase.AddStock(c.Request.Context(), idInt, stockUpdate.Amount); err != nil {
		c.Error(err)
		c.Next()
		return
	}

	response := response.Success("stock updated")
	c.JSON(http.StatusOK, response)
}

// List retrieves a list of goods with optional filtering and pagination
func (h *GoodsHandler) List(c *gin.Context) {
	var query dto.GoodsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		err := custom_error.NewError(custom_error.ErrBadRequest, err.Error())
		c.Error(err)
		c.Next()
		return
	}

	goods, total, err := h.goodsUsecase.List(c.Request.Context(), query)
	if err != nil {
		c.Error(err)
		c.Next()
		return
	}
	data := map[string]interface{}{
		"goods": goods,
		"total": total,
	}
	response := response.Success(data)
	c.JSON(http.StatusOK, response)
}
