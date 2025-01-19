package models

import (
	"errors"
	"fmt"
)

type Item struct {
	ShortDescription string
	Price            string
	priceFloat       float64
}

// IsValid returns an error if the Item object is not valid.
func (item Item) IsValid() (err error) {
	if item.ShortDescription == "" {
		err = errors.Join(err, ErrItemShortDescriptionBlank)
	} else if !shortDescriptionRegex.MatchString(item.ShortDescription) {
		err = errors.Join(err, fmt.Errorf("%s is an %w", item.ShortDescription, ErrItemShortDescriptionInvalid))
	}

	if item.Price == "" {
		err = errors.Join(err, ErrItemPriceBlank)
	} else if !priceRegex.MatchString(item.Price) {
		err = errors.Join(err, fmt.Errorf("%s is an %w", item.Price, ErrPriceFormatInvalid))
	}

	if err != nil {
		err = fmt.Errorf("%w: %w", ErrItemInvalid, err)
	}
	return err
}
