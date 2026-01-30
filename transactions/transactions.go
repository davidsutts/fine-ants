package transactions

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ausocean/cloud/datastore"
)

const TypeTransaction = "Transaction"

type Transaction struct {
	ID            string
	EffectiveDate time.Time
	EnteredDate   time.Time
	Description   string
	Amount        float64
	Balance       float64
	Category      string
}

const CatUncategorised string = "uncategorised"

// Copy is not currently implemented.
func (s *Transaction) Copy(datastore.Entity) (datastore.Entity, error) {
	return nil, datastore.ErrUnimplemented
}

// GetCache returns nil, indicating no caching.
func (s *Transaction) GetCache() datastore.Cache {
	return nil
}

func ParseTransactionsFromCSV(data []byte) ([]Transaction, error) {
	r := csv.NewReader(bytes.NewBuffer(data))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("unable to parse transactions with CSV reader: %v", err)
	}

	trans := []Transaction{}
	for _, row := range records[1:] {

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

		amount, err := strconv.ParseFloat(strings.ReplaceAll(row[3], "$", ""), 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse amount as float: %v", err)
		}

		balance, err := strconv.ParseFloat(strings.ReplaceAll(row[4], "$", ""), 64)
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

func GetTransactionHash(effectiveDate time.Time, amount, balance float64, desc string) string {
	data := fmt.Appendf(nil, "%v:%f:%f:%s", effectiveDate, amount, balance, desc)
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

// Create will try to create a transaction.
func (t *Transaction) Create(ctx context.Context, store datastore.Store) error {
	t.ID = GetTransactionHash(t.EffectiveDate, t.Amount, t.Balance, t.Description)
	if t.Category == "" {
		t.Category = CatUncategorised
	}
	key := store.NameKey(TypeTransaction, t.ID)
	return store.Create(ctx, key, t)
}

func GetAll(ctx context.Context, store datastore.Store) ([]Transaction, error) {
	q := store.NewQuery(TypeTransaction, false, "ID")
	txs := []Transaction{}
	_, err := store.GetAll(ctx, q, &txs)
	if err != nil {
		return nil, fmt.Errorf("unable to get all transactions: %v", err)
	}
	return txs, nil
}
