/*
NAME
  Fine Ants - Finance done properly

AUTHORS
  David Sutton <dsutton1202@gmail.com>

LICENSE
  Copyright (C) 2026 - Echo-dass Software

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/ausocean/cloud/datastore"
	"github.com/davidsutts/fine-ants/transactions"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const version = "0.0.1"

// service holds important information about the app accessible by all
// handlers.
type service struct {
	store datastore.Store
	dev   bool
}

// registerEndpoints registers the endpoints to be used in the app.
func (svc *service) registerEndpoints(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/version", svc.versionHandler)

	api.Group("/transaction").
		Post("/parse", svc.parseTransactionsHandler).
		Post("/categorise", svc.categoriseTransactionHandler).
		Get("/get-all", svc.getAllTransactionsHandler).
		Get("/download", svc.downloadTransactionsHandler)
}

func errorHandler(c *fiber.Ctx, err error) error {
	log.Errorf("[%s] %v", c.Route().Path, err)
	c.Status(fiber.StatusInternalServerError)
	return nil
}

func main() {
	svc := &service{}

	ctx := context.Background()

	// Get CLI args.
	port := flag.Int("port", 8080, "Port number to host webserver")
	dev := flag.Bool("dev", false, "Run in development mode")
	flag.Parse()

	svc.dev = *dev
	if svc.dev {
		log.Info("running in dev mode")
	}

	svc.getStore(ctx)

	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	svc.registerEndpoints(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", *port)))
}

// versionHandler handles requests to /version.
func (svc *service) versionHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"version": version})
}

func (svc *service) parseTransactionsHandler(c *fiber.Ctx) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return fmt.Errorf("unable to get formfile: %v", err)
	}

	file, err := fh.Open()
	if err != nil {
		return fmt.Errorf("unable to open formfile: %v", err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return c.JSON(fiber.Map{"error": fmt.Sprintf("unable to read all transactions: %v", err)})
	}

	txs, err := transactions.ParseTransactionsFromCSV(data)
	if err != nil {
		return c.JSON(fiber.Map{"error": fmt.Sprintf("unable to parse transactions: %v", err)})
	}

	for _, tx := range txs {
		err := tx.Create(c.Context(), svc.store)
		if err != nil {
			log.Errorf("unable to create transaction (%s): %v", tx.Description, err)
		}
	}

	return c.JSON(nil)

}

func (svc *service) getAllTransactionsHandler(c *fiber.Ctx) error {
	txs, err := transactions.GetAll(c.Context(), svc.store)
	if err != nil {
		return err
	}

	slices.SortFunc(txs, func(a, b transactions.Transaction) int {
		if a.EffectiveDate.Before(b.EffectiveDate) {
			return -1
		}
		if b.EffectiveDate.Before(a.EffectiveDate) {
			return 1
		}

		// Both dates are the same, sort by balance.
		if a.Balance+b.Amount == b.Balance {
			return -1
		} else {
			return 1
		}

	})

	return c.JSON(txs)
}

func (svc *service) categoriseTransactionHandler(c *fiber.Ctx) error {
	tx := transactions.Transaction{}
	err := json.Unmarshal(c.Body(), &tx)
	if err != nil {
		return fmt.Errorf("unable to unmarshal transaction from body: %v", err)
	}

	err = tx.UpdateCategory(c.Context(), svc.store, tx.Category)
	if err != nil {
		return fmt.Errorf("unable to update transaction category: %v", err)
	}

	log.Infof("Categorised transaction: %+v", tx)

	return c.JSON(tx)
}

func (svc *service) downloadTransactionsHandler(c *fiber.Ctx) error {
	txs, err := transactions.GetAll(c.Context(), svc.store)
	if err != nil {
		return err
	}

	w := csv.NewWriter(c)
	defer w.Flush()
	for _, tx := range txs {
		err = w.Write([]string{tx.EffectiveDate.String(), strconv.FormatFloat(tx.Amount, 'f', 2, 64), tx.Category})
		if err != nil {
			log.Errorf("failed to write record")
		}
	}

	c.Response().Header.SetContentType("text/csv")
	c.Response().Header.Set("Content-Disposition", "attachment; filename=transactions.csv")

	return nil
}
