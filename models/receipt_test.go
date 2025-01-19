package models

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReceiptIsValid(t *testing.T) {
	testDate, err := time.Parse(time.DateOnly, "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}

	testTime, err := time.Parse(timeFormat, "13:01")
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
				Retailer:     "23456715215!",
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
	testTime1, err := time.Parse(timeFormat, "13:01")
	if err != nil {
		t.Fatal(err)
	}
	testDate2, err := time.Parse(time.DateOnly, "2022-03-20")
	if err != nil {
		t.Fatal(err)
	}
	testTime2, err := time.Parse(timeFormat, "14:33")
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
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49", priceFloat: 6.49},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25", priceFloat: 12.25},
					{ShortDescription: "Knorr Creamy Chicken", Price: "1.26", priceFloat: 1.26},
					{ShortDescription: "Doritos Nacho Cheese", Price: "3.35", priceFloat: 3.35},
					{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00", priceFloat: 12.00},
				},
				Total:      "35.35",
				totalFloat: 35.35,
			},
			wantPoints: 28,
		},
		{
			receipt: Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: testDate2,
				PurchaseTime: testTime2,
				Items: []Item{
					{ShortDescription: "Gatorade", Price: "2.25", priceFloat: 2.25},
					{ShortDescription: "Gatorade", Price: "2.25", priceFloat: 2.25},
					{ShortDescription: "Gatorade", Price: "2.25", priceFloat: 2.25},
					{ShortDescription: "Gatorade", Price: "2.25", priceFloat: 2.25},
				},
				Total:      "9.00",
				totalFloat: 9.00,
			},
			wantPoints: 109,
		},
		{
			receipt: Receipt{
				Retailer:     "test-retailer",
				PurchaseDate: testDate1,
				PurchaseTime: testTime2,
				Items: []Item{
					{ShortDescription: "Something cheap", Price: "0.70", priceFloat: 0.70},
				},
				Total:      "0.70",
				totalFloat: 0.70,
			},
			wantPoints: 12 + 6 + 10 + 1,
		},
		{
			receipt: Receipt{
				Retailer:     "test-retailer",
				PurchaseDate: testDate1,
				PurchaseTime: testTime2,
				Items: []Item{
					{ShortDescription: "Something free", Price: "0.00", priceFloat: 0.00},
				},
				Total:      "0.00",
				totalFloat: 0.00,
			},
			wantPoints: 12 + 6 + 10,
		},
	}

	for _, tc := range testcases {
		points := tc.receipt.CalculatePoints()

		if points != tc.wantPoints {
			t.Errorf("%+v.CalculatePoints(); got: %d, want: %d", tc.receipt, points, tc.wantPoints)
		}
	}
}

func TestReceiptUnmarshalJSON(t *testing.T) {

	testDate1, err := time.Parse(time.DateOnly, "2022-01-02")
	if err != nil {
		t.Fatal(err)
	}
	testTime1, err := time.Parse(timeFormat, "08:13")
	if err != nil {
		t.Fatal(err)
	}
	testDate2, err := time.Parse(time.DateOnly, "2022-01-02")
	if err != nil {
		t.Fatal(err)
	}
	testTime2, err := time.Parse(timeFormat, "13:13")
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		input   string
		want    Receipt
		wantErr error
	}{
		{
			input: `{
						"retailer": "Walgreens",
						"purchaseDate": "2022-01-02",
						"purchaseTime": "08:13",
						"total": "2.65",
						"items": [
							{"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
							{"shortDescription": "Dasani", "price": "1.40"}
						]
					}`,
			want: Receipt{
				Retailer:     "Walgreens",
				PurchaseDate: testDate1,
				PurchaseTime: testTime1,
				Total:        "2.65",
				totalFloat:   2.65,
				Items: []Item{
					{ShortDescription: "Pepsi - 12-oz", Price: "1.25", priceFloat: 1.25},
					{ShortDescription: "Dasani", Price: "1.40", priceFloat: 1.40},
				},
			},
			wantErr: nil,
		},
		{
			input: `{
						"retailer": "Target",
						"purchaseDate": "2022-01-02",
						"purchaseTime": "13:13",
						"total": "1.25",
						"items": [
							{"shortDescription": "Pepsi - 12-oz", "price": "1.25"}
						]
					}`,
			want: Receipt{
				Retailer:     "Target",
				PurchaseDate: testDate2,
				PurchaseTime: testTime2,
				Total:        "1.25",
				totalFloat:   1.25,
				Items: []Item{
					{ShortDescription: "Pepsi - 12-oz", Price: "1.25", priceFloat: 1.25},
				},
			},
			wantErr: nil,
		},
		{
			input:   `{"retailer": "", "purchaseDate": "", "purchaseTime":"","total":"","items":[]}`,
			wantErr: ErrReceiptInvalid,
		},
	}

	for _, tc := range testcases {
		var testReceipt Receipt
		if err := json.Unmarshal([]byte(tc.input), &testReceipt); err != nil {
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Unmarshal(%s) returned an unexpected error: %v", tc.input, err)
			}
			continue
		}
		if tc.wantErr != nil {
			t.Errorf("Unmarshal(%s) expected error: %v", tc.input, tc.wantErr)
			continue
		}

		assert.Equal(t, tc.want, testReceipt)
	}
}
