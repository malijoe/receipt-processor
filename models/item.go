package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Item struct {
	ShortDescription string
	Price            string
	priceFloat       float64
}

// IsValid returns an error if the Item object is not valid.
func (item *Item) IsValid() (err error) {
	if item.ShortDescription == "" {
		err = errors.Join(err, ErrItemShortDescriptionBlank)
	} else if !shortDescriptionRegex.MatchString(item.ShortDescription) {
		err = errors.Join(err, fmt.Errorf("%s is an %w", item.ShortDescription, ErrItemShortDescriptionInvalid))
	}

	if item.Price == "" {
		err = errors.Join(err, ErrItemPriceBlank)
	} else if !priceRegex.MatchString(item.Price) {
		err = errors.Join(err, fmt.Errorf("%s is an %w", item.Price, ErrPriceFormatInvalid))
	} else {
		// if a price is provided, parse the float. otherwise let the validation method catch the error.
		priceFloat, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return err
		}
		item.priceFloat = priceFloat
	}

	if err != nil {
		err = fmt.Errorf("%w: %w", ErrItemInvalid, err)
	}
	return err
}

// Unmarshal handles generic unmarshalling for item object
func (item *Item) Unmarshal(unmarshal func(any) error) error {
	var obj struct {
		ShortDescription string `json:"shortDescription"`
		Price            string `json:"price"`
	}

	if err := unmarshal(&obj); err != nil {
		return err
	}

	item.ShortDescription = obj.ShortDescription
	item.Price = obj.Price

	return nil
}

// UnmarshalJSON handles unmarshalling JSON data into an Item object.
func (item *Item) UnmarshalJSON(data []byte) error {
	return item.Unmarshal(func(obj any) error {
		return json.Unmarshal(data, obj)
	})
}
