package models

import (
	"errors"
	"testing"
	"time"
)

func TestReceiptIsValid(t *testing.T) {
	testDate, err := time.Parse(time.DateOnly, "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}

	testTime, err := time.Parse("15:04", "13:01")
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		receipt  Receipt
		wantErrs []error
	}{
		{
			receipt: Receipt{
				Retailer:     "",
				PurchaseDate: time.Time{},
				PurchaseTime: time.Time{},
				Items:        nil,
				Total:        "",
			},
			wantErrs: []error{ErrReceiptInvalid, ErrReceiptRetailerBlank, ErrReceiptPurchaseDateBlank, ErrReceiptPurchaseTimeBlank, ErrReceiptTotalBlank},
		},
		{
			receipt: Receipt{
				Retailer:     "apple",
				PurchaseDate: testDate,
				PurchaseTime: testTime,
				Items:        []Item{{ShortDescription: "test-item"}},
				Total:        "00000",
			},
			wantErrs: []error{ErrReceiptInvalid, ErrReceiptRetailerInvalid, ErrPriceFormatInvalid},
		},
		{
			receipt: Receipt{
				Retailer:     `w-s&`,
				PurchaseDate: testDate,
				PurchaseTime: testTime,
				Items:        []Item{{ShortDescription: "test-item"}},
				Total:        "42.00",
			},
		},
	}

	for _, tc := range testcases {
		if err := tc.receipt.IsValid(); err != nil {
			if len(tc.wantErrs) == 0 {
				t.Errorf("%v.IsValid(); found unexpected error: %v", tc.receipt, err)
			}
			for _, wantErr := range tc.wantErrs {
				if !errors.Is(err, wantErr) {
					t.Errorf("%v.IsValid(); did not find expected error %v", tc.receipt, wantErr)
				}
			}
			continue
		}
		if len(tc.wantErrs) > 0 {
			t.Errorf("%v.IsValid(); expected error(s) %v, but none were thrown", tc.receipt, tc.wantErrs)
		}
	}
}

func TestReceiptCalculatePoints(t *testing.T) {
	testDate1, err := time.Parse(time.DateOnly, "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}
	testTime1, err := time.Parse("15:04", "13:01")
	if err != nil {
		t.Fatal(err)
	}
	testDate2, err := time.Parse(time.DateOnly, "2022-03-20")
	if err != nil {
		t.Fatal(err)
	}
	testTime2, err := time.Parse("15:04", "14:33")
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		receipt    Receipt
		wantPoints int
	}{
		{
			receipt: Receipt{
				Retailer:     "Target",
				PurchaseDate: testDate1,
				PurchaseTime: testTime1,
				Items: []Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
					{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
					{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
					{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
				},
				Total: "35.35",
			},
			wantPoints: 28,
		},
		{
			receipt: Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: testDate2,
				PurchaseTime: testTime2,
				Items: []Item{
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
		points, err := tc.receipt.CalculatePoints()
		if err != nil {
			t.Errorf("%v.CalculatePoints() returned an unexpected error: %v", tc.receipt, err)
			continue
		}

		if points != tc.wantPoints {
			t.Errorf("%v.CalculatePoints(); got: %d, want: %d", tc.receipt, points, tc.wantPoints)
		}
	}
}
