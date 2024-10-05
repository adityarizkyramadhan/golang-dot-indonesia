package dto

type CreateInvoiceRequest struct {
	InvoiceNumber *string         `json:"invoice_number" binding:"required"`
	IsPaid        *bool           `json:"is_paid"`
	Purchases     []PurchaseInput `json:"purchases" binding:"required,dive"`
}

type PurchaseInput struct {
	GoodsID      *int     `json:"goods_id" binding:"required"`
	Amount       *int     `json:"amount" binding:"required"`
	PricePerUnit *float64 `json:"price_per_unit" binding:"required"`
}
