package models

import (
	"errors"
	"regexp"
)

var (
	// error stubs for receipt object
	ErrReceiptRetailerBlank     = errors.New("receipt retailer cannot be blank")
	ErrReceiptRetailerInvalid   = errors.New("invalid receipt retailer")
	ErrReceiptPurchaseDateBlank = errors.New("receipt purchase date cannot be blank")
	ErrReceiptPurchaseTimeBlank = errors.New("receipt purchase time cannot be blank")
	ErrReceiptItemsEmpty        = errors.New("receipt must have items")
	ErrReceiptTotalBlank        = errors.New("receipt total cannot be blank")
	ErrReceiptInvalid           = errors.New("invalid receipt")

	// error stubs for item object
	ErrItemShortDescriptionBlank   = errors.New("item short description cannot be blank")
	ErrItemShortDescriptionInvalid = errors.New("invalid item short description")
	ErrItemPriceBlank              = errors.New("item price cannot be blank")
	ErrItemInvalid                 = errors.New("invalid item")

	// general error stubs
	ErrPriceFormatInvalid = errors.New("invalid price format")

	// regex to validate retailer field
	retailerRegex = regexp.MustCompile(`^[\w\s\-&]+$`)
	// regex to validate price formated strings.
	priceRegex = regexp.MustCompile(`^\d+\.\d{2}$`)
	// regex to validate the shortDescription field for Item objects.
	shortDescriptionRegex = regexp.MustCompile(`^[\w\s\-]+$`)
	// regex for catching all individual alphanumeric characters.
	alphanumericRegex = regexp.MustCompile(`[\w\d]`)
)
