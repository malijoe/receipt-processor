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
	retailerRegex = regexp.MustCompile(`^[\w\s\-&]+$`)
	// regex to validate price formated strings.
	priceRegex = regexp.MustCompile(`^\d+\.\d{2}$`)
	// regex to validate the shortDescription field for Item objects.
	shortDescriptionRegex = regexp.MustCompile(`^[\w\s\-]+$`)
	// regex for catching all individual alphanumeric characters.
	alphanumericRegex = regexp.MustCompile(`[\w\d]`)
)
