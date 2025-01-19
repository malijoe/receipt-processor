package models

import (
	"errors"
	"testing"
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
