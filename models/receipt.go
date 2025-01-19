package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "15:04"

type Receipt struct {
	Retailer     string
	PurchaseDate time.Time
	PurchaseTime time.Time
	Items        []Item
	Total        string
	totalFloat   float64
}

func (r Receipt) IsValid() (err error) {
	if r.Retailer == "" {
		err = errors.Join(err, ErrReceiptRetailerBlank)
	} else if !retailerRegex.MatchString(r.Retailer) {
		err = errors.Join(err, fmt.Errorf("%s is an %w", r.Retailer, ErrReceiptRetailerInvalid))
	}

	if r.PurchaseDate.IsZero() {
		err = errors.Join(err, ErrReceiptPurchaseDateBlank)
	}

	if r.PurchaseTime.IsZero() {
		err = errors.Join(err, ErrReceiptPurchaseTimeBlank)
	}

	if len(r.Items) < 1 {
		err = errors.Join(err, ErrReceiptItemsEmpty)
	}

	if r.Total == "" {
		err = errors.Join(err, ErrReceiptTotalBlank)
	} else if !priceRegex.MatchString(r.Total) {
		err = errors.Join(err, fmt.Errorf("%s is an %w", r.Total, ErrPriceFormatInvalid))
	}

	if err != nil {
		err = fmt.Errorf("%w: %w", ErrReceiptInvalid, err)
	}
	return err
}

// CalculatePoints returns the number of points earned by the receipt.
func (r Receipt) CalculatePoints() (points int) {
	// count the number of alphanumeric characters and add that to the number of points
	points += len(alphanumericRegex.FindAllString(r.Retailer, -1))
	// determine whether the total amount is a whole number
	isWhole := math.Ceil(r.totalFloat) == r.totalFloat
	if isWhole && r.totalFloat > 1 {
		// add 50 pts if the total is a round dollar amount with no cents
		points += 50
	}

	// get just the dollar amount
	dollars := math.Floor(r.totalFloat)
	// subtract the dollar amount from the total to get the cents.
	cents := r.totalFloat - dollars
	// multiple the cents by 100 to get whole numbers and cast to integer to avoid float math
	adjustedCents := int(cents * 100)
	if adjustedCents%25 == 0 && r.totalFloat > 1 {
		// add 25 pts if the quantity of cents is a multiple of 0.25
		points += 25
	}

	numItems := len(r.Items)
	// add 5 points for every two items on the receipt
	pointsFromNumItems := 5 * (numItems / 2)
	points += pointsFromNumItems

	for _, item := range r.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			// when the trimmed length of the item description is a multiple of 3
			// multiple the price by 0.2 and round up to the nearest integer
			price := item.priceFloat * 0.2

			// add .5 to price before rounding so that we will always round up, then add the points
			pointsFromItem := int(math.Round(price + 0.5))
			points += pointsFromItem
		}
	}

	if r.PurchaseDate.Day()%2 == 1 {
		// add 6 pts if purchase day is odd
		points += 6
	}
	hour := r.PurchaseTime.Hour()
	minute := r.PurchaseTime.Minute()

	if ((hour == 14 && minute > 0) || hour > 14) && hour < 16 {
		// add 10 pts if the time of purchase is after 2pm and before 4pm
		points += 10
	}

	return points
}

// Unmarshal handles generic unmarshalling for the receipt object.
func (r *Receipt) Unmarshal(unmarshal func(any) error) error {
	var obj struct {
		Retailer     string `json:"retailer"`
		PurchaseDate string `json:"purchaseDate"`
		PurchaseTime string `json:"purchaseTime"`
		Total        string `json:"total"`
		Items        []Item `json:"items"`
	}

	if err := unmarshal(&obj); err != nil {
		return err
	}

	if obj.PurchaseDate != "" {
		// if a purchase date is provided, parse it. otherwise, let the validation method catch the error.
		purchaseDate, err := time.Parse(time.DateOnly, obj.PurchaseDate)
		if err != nil {
			return err
		}
		r.PurchaseDate = purchaseDate
	}

	if obj.PurchaseTime != "" {
		// if a purchase time is provided, parse it. otherwise, let the validation method catch the error.
		purchaseTime, err := time.Parse(timeFormat, obj.PurchaseTime)
		if err != nil {
			return err
		}
		r.PurchaseTime = purchaseTime
	}

	r.Retailer = obj.Retailer
	r.Total = obj.Total
	r.Items = obj.Items
	if obj.Total != "" {
		// if a total is provided, parse the float. otherwise, let the validation method catch the error.
		totalFloat, err := strconv.ParseFloat(obj.Total, 64)
		if err != nil {
			return err
		}
		r.totalFloat = totalFloat
	}

	if err := r.IsValid(); err != nil {
		return err
	}
	return nil
}

// UnmarshalJSON handles unmarshalling JSON data into a Receipt object.
func (r *Receipt) UnmarshalJSON(data []byte) error {
	return r.Unmarshal(func(obj any) error {
		return json.Unmarshal(data, obj)
	})
}
