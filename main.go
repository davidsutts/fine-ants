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
	"flag"
	"fmt"
	"io"

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
		Get("/get-all", svc.getAllTransactionsHandler)
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

	transactions, err := transactions.ParseTransactionsFromCSV(data)
	if err != nil {
		return c.JSON(fiber.Map{"error": fmt.Sprintf("unable to parse transactions: %v", err)})
	}

	return c.JSON(transactions)

}

func (svc *service) getAllTransactionsHandler(c *fiber.Ctx) error {
	txs, err := transactions.GetAll(c.Context(), svc.store)
	if err != nil {
		return err
	}

	return c.JSON(txs)
}
