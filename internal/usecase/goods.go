package usecase

import (
	"context"
	"strconv"

	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/dto"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/entity"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/repository"
	custom_error "github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/errors"
)

type Goods struct {
	goodRepo *repository.Good
}

type GoodsUsecase interface {
	Create(ctx context.Context, good *dto.Goods) error
	Get(ctx context.Context, id int) (*entity.Goods, error)
	Update(ctx context.Context, good *dto.GoodsUpdate, id *string) error
	Delete(ctx context.Context, id int) error
	AddStock(ctx context.Context, id int, amount int) error
	List(ctx context.Context, query dto.GoodsQuery) ([]entity.Goods, int64, error)
}

func NewGoods(goodRepo *repository.Good) GoodsUsecase {
	return &Goods{goodRepo: goodRepo}
}

// Create adds a new Good
func (uc *Goods) Create(ctx context.Context, good *dto.Goods) error {
	if good.Name == nil {
		return custom_error.NewError(custom_error.ErrBadRequest, "good name cannot be nil")
	}
	if good.Price == nil {
		return custom_error.NewError(custom_error.ErrBadRequest, "good price cannot be nil")
	}
	if good.Stock == nil {
		return custom_error.NewError(custom_error.ErrBadRequest, "good stock cannot be nil")
	}
	goods := &entity.Goods{
		Name:  good.Name,
		Price: good.Price,
		Stock: good.Stock,
	}
	return uc.goodRepo.Create(ctx, nil, goods)
}

// Get retrieves a Good by ID
func (uc *Goods) Get(ctx context.Context, id int) (*entity.Goods, error) {
	return uc.goodRepo.Get(ctx, nil, id)
}

// Update modifies an existing Good
func (uc *Goods) Update(ctx context.Context, good *dto.GoodsUpdate, id *string) error {
	if id == nil {
		return custom_error.NewError(custom_error.ErrBadRequest, "id cannot be nil")
	}

	idInt, err := strconv.Atoi(*id)
	if err != nil {
		return custom_error.NewError(custom_error.ErrBadRequest, "invalid ID")
	}

	goodsData, err := uc.goodRepo.Get(ctx, nil, idInt)
	if err != nil {
		return custom_error.NewError(custom_error.ErrNotFound, "good not found")
	}

	if good.Name != nil {
		goodsData.Name = good.Name
	}

	if good.Price != nil {
		goodsData.Price = good.Price
	}

	if good.Stock != nil {
		goodsData.Stock = good.Stock
	}

	return uc.goodRepo.Update(ctx, nil, goodsData)
}

// Delete removes a Good by ID
func (uc *Goods) Delete(ctx context.Context, id int) error {
	return uc.goodRepo.Delete(ctx, nil, id)
}

// AddStock increases the stock of a Good
func (uc *Goods) AddStock(ctx context.Context, id int, amount int) error {
	if amount <= 0 {
		return custom_error.NewError(custom_error.ErrBadRequest, "amount must be greater than 0")
	}
	return uc.goodRepo.AddStock(ctx, nil, id, amount)
}

// List retrieves a list of Goods with optional filtering and pagination
func (uc *Goods) List(ctx context.Context, query dto.GoodsQuery) ([]entity.Goods, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	return uc.goodRepo.List(ctx, query.Name, query.Page, query.PageSize)
}
