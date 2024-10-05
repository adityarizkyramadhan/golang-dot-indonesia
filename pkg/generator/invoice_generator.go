package generator

import (
	"fmt"
	"math/rand"
	"time"
)

func generateInvoiceNumber() string {
	return fmt.Sprintf("INV-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000000))
}

func InvoiceGenerator() *string {
	invoiceNumber := generateInvoiceNumber()
	return &invoiceNumber
}
