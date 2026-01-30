package main

import (
	"context"

	"github.com/ausocean/cloud/datastore"
	"github.com/davidsutts/fine-ants/transactions"
	"github.com/gofiber/fiber/v2/log"
)

func registerEntities(store datastore.Store) {
	datastore.RegisterEntity(transactions.TypeTransaction, func() datastore.Entity { return &transactions.Transaction{} })
}

func (svc *service) getStore(ctx context.Context) {
	kind := "cloud"
	if svc.dev {
		kind = "file"
	}
	var err error
	svc.store, err = datastore.NewStore(ctx, kind, "fine-ants", "")
	if err != nil {
		log.Fatalf("unable to get store: %v", err)
	}

	registerEntities(svc.store)

	log.Infof("created new %sstore", kind)
}
