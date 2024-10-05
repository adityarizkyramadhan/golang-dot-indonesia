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

type InvoiceHandler struct {
	invoiceUsecase *usecase.InvoiceUsecase
}

func NewInvoiceHandler(invoiceUsecase *usecase.InvoiceUsecase) *InvoiceHandler {
	return &InvoiceHandler{invoiceUsecase: invoiceUsecase}
}

func (h *InvoiceHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/", middleware.JWTMiddleware(), h.CreateInvoice)
	r.GET("/:id", middleware.JWTMiddleware(), h.GetInvoice)
	r.PUT("/:id", middleware.JWTMiddleware(), h.UpdateInvoice)
	r.DELETE("/:id", middleware.JWTMiddleware(), h.DeleteInvoice)
	r.GET("/", middleware.JWTMiddleware(), h.GetAllInvoices)
}

// CreateInvoice handles the request for creating an invoice
func (h *InvoiceHandler) CreateInvoice(ctx *gin.Context) {
	var request dto.CreateInvoiceRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errResponse := custom_error.NewError(custom_error.ErrBadRequest, err.Error())
		ctx.Error(errResponse)
		ctx.Next()
		return
	}

	userID := ctx.MustGet("id").(int)

	request.UserID = &userID

	if err := h.invoiceUsecase.CreateInvoice(ctx.Request.Context(), request); err != nil {
		ctx.Error(err)
		ctx.Next()
		return
	}

	response := response.Success("Invoice created")
	ctx.JSON(http.StatusCreated, response)
}

// GetInvoice handles the request for retrieving an invoice by ID
func (h *InvoiceHandler) GetInvoice(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errResponse := custom_error.NewError(custom_error.ErrBadRequest, "invalid ID")
		ctx.Error(errResponse)
		ctx.Next()
		return
	}

	invoice, err := h.invoiceUsecase.GetInvoice(ctx.Request.Context(), idInt, &userID)
	if err != nil {
		ctx.Error(err)
		ctx.Next()
		return
	}

	response := response.Success(invoice)
	ctx.JSON(http.StatusOK, response)
}

// UpdateInvoice handles the request for updating an existing invoice
func (h *InvoiceHandler) UpdateInvoice(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)
	var invoice entity.InvoicePurchase
	if err := ctx.ShouldBindJSON(&invoice); err != nil {
		errResponse := custom_error.NewError(custom_error.ErrBadRequest, err.Error())
		ctx.Error(errResponse)
		ctx.Next()
		return
	}

	if err := h.invoiceUsecase.UpdateInvoice(ctx.Request.Context(), &invoice, &userID); err != nil {
		ctx.Error(err)
		ctx.Next()
		return
	}

	response := response.Success("Invoice updated")
	ctx.JSON(http.StatusOK, response)
}

// DeleteInvoice handles the request for deleting an invoice by ID
func (h *InvoiceHandler) DeleteInvoice(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errResponse := custom_error.NewError(custom_error.ErrBadRequest, "invalid ID")
		ctx.Error(errResponse)
		ctx.Next()
		return
	}

	if err := h.invoiceUsecase.DeleteInvoice(ctx.Request.Context(), idInt, &userID); err != nil {
		ctx.Error(err)
		ctx.Next()
		return
	}

	response := gin.H{"message": "Invoice deleted successfully"}
	ctx.JSON(http.StatusOK, response)
}

// GetAllInvoices handles the request for retrieving all invoices
func (h *InvoiceHandler) GetAllInvoices(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)
	invoices, err := h.invoiceUsecase.GetAllInvoices(ctx.Request.Context(), &userID)
	if err != nil {
		ctx.Error(err)
		ctx.Next()
		return
	}

	response := response.Success(invoices)
	ctx.JSON(http.StatusOK, response)
}
