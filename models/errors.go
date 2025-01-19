package models

import (
	"errors"
	"regexp"
)

var (
	ErrReceiptRetailerBlank     = errors.New("receipt retailer cannot be blank")
	ErrReceiptRetailerInvalid   = errors.New("invalid receipt retailer")
	ErrReceiptPurchaseDateBlank = errors.New("receipt purchase date cannot be blank")
	ErrReceiptPurchaseTimeBlank = errors.New("receipt purchase time cannot be blank")
	ErrReceiptItemsEmpty        = errors.New("receipt must have items")
	ErrReceiptTotalBlank        = errors.New("receipt total cannot be blank")
	ErrPriceFormatInvalid       = errors.New("invalid price format")
	ErrReceiptInvalid           = errors.New("invalid receipt")

	// regex to validate retailer field
	// provided pattern for retailer field '^[\\w\\s\\-&]+$' did not compile so the pattern was changed
	// to match the one provided for the shortDescription field of the item object.
	retailerRegex = regexp.MustCompile(`^[\\w\\s\\-]+$`)
	// regex to validate price formated strings.
	priceRegex = regexp.MustCompile(`^\d+\.\d{2}$`)
	// regex to validate the shortDescription field for Item objects.
	// Accepts only the following characters: "w","s","-","\".
	shortDescriptionRegex = regexp.MustCompile(`^[\\w\\s\\-]+$`)
	// regex for catching all individual alphanumeric characters.
	alphanumericRegex = regexp.MustCompile(`[\w\d]`)
)
