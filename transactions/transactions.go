package transactions

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	EffectiveDate time.Time
	EnteredDate   time.Time
	Description   string
	Amount        float64
	Balance       float64
}

func ParseTransactionsFromCSV(data []byte) ([]Transaction, error) {
	r := csv.NewReader(bytes.NewBuffer(data))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("unable to parse transactions with CSV reader: %v", err)
	}

	trans := []Transaction{}
	for _, row := range records[1:] {
		log.Printf("row: %+v", row)

		enteredDate, err := time.Parse("02/01/2006", row[1])
		if err != nil {
			return nil, fmt.Errorf("unable to parse entered date: %v", err)
		}

		var effectiveDate time.Time
		if row[0] == "" {
			effectiveDate = enteredDate
		} else {
			effectiveDate, err = time.Parse("02/01/2006", row[0])
			if err != nil {
				return nil, fmt.Errorf("unable to parse effective date: %v", err)
			}
		}

		amount, err := strconv.ParseFloat(strings.Split(row[3], "$")[1], 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse amount as float: %v", err)
		}

		balance, err := strconv.ParseFloat(strings.Split(row[4], "$")[1], 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse balance as float: %v", err)
		}

		t := Transaction{
			EffectiveDate: effectiveDate,
			EnteredDate:   enteredDate,
			Description:   row[2],
			Amount:        amount,
			Balance:       balance,
		}

		trans = append(trans, t)
	}

	return trans, nil
}
