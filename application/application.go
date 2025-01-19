package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/malijoe/receipt-processor/models"
	statuserrors "github.com/malijoe/receipt-processor/statusErrors"
)

type Application struct {
	data map[string]*models.Receipt
}

func NewApplication() *Application {
	return &Application{
		data: make(map[string]*models.Receipt),
	}
}

// ProcessReceipt takes a receipt object saves it to memory and returns the generated id for the receipt.
func (app *Application) ProcessReceipt(ctx context.Context, receipt models.Receipt) (id string, _ error) {
	// make sure the passed receipt is valid
	if err := receipt.IsValid(); err != nil {
		return "", fmt.Errorf("%w: %w", statuserrors.ErrBadRequest, err)
	}

	id = uuid.NewString()
	app.data[id] = &receipt
	return id, nil
}

func (app *Application) GetReceiptPoints(ctx context.Context, receiptId string) (points int, _ error) {
	receipt, hasReceipt := app.data[receiptId]
	if !hasReceipt {
		return 0, fmt.Errorf("%w: no receipt found with id %s", statuserrors.ErrNotFound, receiptId)
	}

	pts := receipt.CalculatePoints()
	return pts, nil
}
