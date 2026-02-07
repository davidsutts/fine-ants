package accounts

import (
	"context"
	"fmt"

	"github.com/ausocean/cloud/datastore"
	"github.com/google/uuid"
)

const TypeAccount = "Account"

type Account struct {
	UUID string
	Name string
}

// Copy is not currently implemented.
func (a *Account) Copy(datastore.Entity) (datastore.Entity, error) {
	return nil, datastore.ErrUnimplemented
}

// GetCache returns nil, indicating no caching.
func (a *Account) GetCache() datastore.Cache {
	return nil
}

// CreateAccount creates a new account with a name, and generates a new UUID.
func CreateAccount(ctx context.Context, store datastore.Store, name string) error {
	a := Account{
		UUID: uuid.NewString(),
		Name: name,
	}
	key := store.NameKey(TypeAccount, a.UUID)
	return store.Create(ctx, key, &a)
}

func UpdateAccount(ctx context.Context, store datastore.Store, uuid, name string) error {
	key := store.NameKey(TypeAccount, uuid)
	acc := Account{}
	return store.Update(ctx, key, func(e datastore.Entity) {
		acc_, ok := e.(*Account)
		if !ok {
			return
		}
		acc_.Name = name
	}, &acc)
}

func GetAllAccounts(ctx context.Context, store datastore.Store) ([]Account, error) {
	accs := []Account{}
	q := store.NewQuery(TypeAccount, false, "UUID")
	_, err := store.GetAll(ctx, q, &accs)
	if err != nil {
		return nil, fmt.Errorf("unable to get all accounts: %v", err)
	}

	return accs, nil
}
