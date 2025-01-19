package application

import (
	"context"
	"testing"
	"time"

	"github.com/malijoe/receipt-processor/models"
	"github.com/stretchr/testify/assert"
)

func TestApplication(t *testing.T) {

	testApp := NewApplication()

	testDate2, err := time.Parse(time.DateOnly, "2022-03-20")
	if err != nil {
		t.Fatal(err)
	}
	testTime2, err := time.Parse("15:04", "14:33")
	if err != nil {
		t.Fatal(err)
	}
	testcases := []struct {
		receipt    models.Receipt
		wantErr    error
		wantPoints int
	}{
		{
			receipt: models.Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: testDate2,
				PurchaseTime: testTime2,
				Items: []models.Item{
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
				},
				Total: "9.00",
			},
			wantPoints: 109,
		},
	}

	for _, tc := range testcases {
		id, err := testApp.ProcessReceipt(context.TODO(), tc.receipt)
		if err != nil {
			t.Errorf("ProcessReceipt(%+v) returned an unexpected error: %v", tc.receipt, err)
			continue
		}

		points, err := testApp.GetReceiptPoints(context.TODO(), id)
		if err != nil {
			t.Errorf("GetReceiptPoints(%s) returned an unexpected error: %v", id, err)
			continue
		}

		assert.Equal(t, tc.wantPoints, points)
	}
}
