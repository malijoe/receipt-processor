package models

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemIsValid(t *testing.T) {
	testcases := []struct {
		item     Item
		wantErrs []error
	}{
		{
			item:     Item{ShortDescription: "", Price: ""},
			wantErrs: []error{ErrItemInvalid, ErrItemShortDescriptionBlank, ErrItemPriceBlank},
		},
		{
			item:     Item{ShortDescription: "&&-wasfdn", Price: "0000000"},
			wantErrs: []error{ErrItemInvalid, ErrItemShortDescriptionInvalid, ErrPriceFormatInvalid},
		},
		{
			item:     Item{ShortDescription: "this-is-a-test", Price: "42.00"},
			wantErrs: nil,
		},
	}
	for _, tc := range testcases {
		if err := tc.item.IsValid(); err != nil {
			if len(tc.wantErrs) == 0 {
				t.Errorf("%v.IsValid(); found an unexpected error: %v", tc.item, err)
			}
			for _, wantErr := range tc.wantErrs {
				if !errors.Is(err, wantErr) {
					t.Errorf("%v.IsValid(); did not find expected error %v", tc.item, wantErr)
				}
			}
			continue
		}
		if len(tc.wantErrs) > 0 {
			t.Errorf("%v.IsValid(); expected error(s) %v, but none were thrown", tc.item, tc.wantErrs)
		}
	}
}

func TestItemUnmarshalJSON(t *testing.T) {
	testcases := []struct {
		input   string
		want    Item
		wantErr error
	}{
		{
			input:   `{"shortDescription": "Pepsi - 12-oz", "price": "1.25"}`,
			want:    Item{ShortDescription: "Pepsi - 12-oz", Price: "1.25", priceFloat: 1.25},
			wantErr: nil,
		},
		{
			input:   `{"shortDescription": "Dasani", "price": "1.40"}`,
			want:    Item{ShortDescription: "Dasani", Price: "1.40", priceFloat: 1.40},
			wantErr: nil,
		},
		{
			input:   `{"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ", "price": "12.00"}`,
			want:    Item{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00", priceFloat: 12.00},
			wantErr: nil,
		},
		{
			input:   `{"shortDescription": "", "price": ""}`,
			wantErr: ErrItemInvalid,
		},
	}

	for _, tc := range testcases {
		var testItem Item
		if err := json.Unmarshal([]byte(tc.input), &testItem); err != nil {
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Unmarshal(%s) returned an unexpected error: %v", tc.input, err)
			}
			continue
		}

		if tc.wantErr != nil {
			t.Errorf("Unmarshal(%s) expected error: %v", tc.input, tc.wantErr)
			continue
		}

		assert.Equal(t, tc.want, testItem)
	}
}
