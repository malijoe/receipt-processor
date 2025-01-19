package models

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type Receipt struct {
	Retailer     string
	PurchaseDate time.Time
	PurchaseTime time.Time
	Items        []Item
	Total        string
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
func (r Receipt) CalculatePoints() (points int, _ error) {
	// count the number of alphanumeric characters and add that to the number of points
	points += len(alphanumericRegex.FindAllString(r.Retailer, -1))
	// separate the whole dollars and cents in the total
	totalPieces := strings.Split(r.Total, ".")
	// parse the pieces of the total to get integers
	totalDollars, err := strconv.Atoi(totalPieces[0])
	if err != nil {
		return 0, err
	}

	totalCents, err := strconv.Atoi(totalPieces[1])
	if err != nil {
		return 0, err
	}

	if totalCents == 0 {
		// add 50 pts if the total is a round dollar amount with no cents
		points += 50
	}
	if totalCents%25 == 0 {
		// add 25 pts if the total is a multiple of 0.25
		points += 25
	} else if totalDollars > 0 && totalCents == 0 {
		// add 25 pts if the total is greater than a dollar and there are no cents.
		// this is a corner case for multiples of 0.25
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
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}

			price = price * 0.2
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

	return points, nil
}
