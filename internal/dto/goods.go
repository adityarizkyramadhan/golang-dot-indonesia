package dto

type Goods struct {
	Name  *string  `json:"name" binding:"required,min=1,max=50"`
	Price *float64 `json:"price" binding:"required,min=1"`
	Stock *int     `json:"stock" binding:"required,min=1"`
}

type GoodsQuery struct {
	Name     *string `form:"name"`
	Page     int     `form:"page"`
	PageSize int     `form:"page_size"`
}
