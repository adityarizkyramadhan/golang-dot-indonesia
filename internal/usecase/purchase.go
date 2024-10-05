package usecase

import (
	"context"

	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/dto"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/entity"
	"github.com/adityarizkyramadhan/golang-dot-indonesia/internal/repository"
	custom_error "github.com/adityarizkyramadhan/golang-dot-indonesia/pkg/errors"
)

type InvoiceUsecase struct {
	invoiceRepo repository.PurchaseRepository
}

func NewInvoiceUsecase(invoiceRepo repository.PurchaseRepository) *InvoiceUsecase {
	return &InvoiceUsecase{invoiceRepo: invoiceRepo}
}

func (uc *InvoiceUsecase) CreateInvoice(ctx context.Context, request dto.CreateInvoiceRequest) error {
	invoicePurchase := &entity.InvoicePurchase{
		InvoiceNumber: request.InvoiceNumber,
		IsPaid:        request.IsPaid,
		UserID:        request.UserID,
	}

	purchases := make([]entity.Purchase, len(request.Purchases))
	for i, item := range request.Purchases {
		purchases[i] = entity.Purchase{
			GoodsID:      item.GoodsID,
			Amount:       item.Amount,
			PricePerUnit: item.PricePerUnit,
			TotalPrice:   calculateTotalPrice(item.Amount, item.PricePerUnit),
		}
	}

	err := uc.invoiceRepo.Create(ctx, invoicePurchase, purchases)
	if err != nil {
		return err
	}

	return nil
}

func (uc *InvoiceUsecase) GetInvoice(ctx context.Context, id int, userID *int) (*entity.InvoicePurchase, error) {
	invoice, err := uc.invoiceRepo.Get(ctx, nil, id, userID)
	if err != nil {
		return nil, custom_error.NewError(custom_error.ErrNotFound, "Invoice not found")
	}
	return invoice, nil
}

func (uc *InvoiceUsecase) UpdateInvoice(ctx context.Context, invoicePurchase *entity.InvoicePurchase, userID *int) error {
	err := uc.invoiceRepo.Update(ctx, nil, invoicePurchase, userID)
	if err != nil {
		return err
	}
	return nil
}

func (uc *InvoiceUsecase) DeleteInvoice(ctx context.Context, id int, userID *int) error {
	err := uc.invoiceRepo.Delete(ctx, nil, id, userID)
	if err != nil {
		return custom_error.NewError(custom_error.ErrNotFound, "Invoice not found")
	}
	return nil
}

func (uc *InvoiceUsecase) GetAllInvoices(ctx context.Context, userID *int) ([]entity.InvoicePurchase, error) {
	invoices, err := uc.invoiceRepo.GetAll(ctx, nil, userID)
	if err != nil {
		return nil, custom_error.NewError(custom_error.ErrNotFound, "Invoices not found")
	}
	return invoices, nil
}

func calculateTotalPrice(amount *int, pricePerUnit *float64) *float64 {
	if amount == nil || pricePerUnit == nil {
		return nil
	}
	totalPrice := float64(*amount) * *pricePerUnit
	return &totalPrice
}
