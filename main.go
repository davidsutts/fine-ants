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
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/davidsutts/fine-ants/transactions"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const version = "0.0.1"

// service holds important information about the app accessible by all
// handlers.
type service struct {
}

// registerEndpoints registers the endpoints to be used in the app.
func (svc *service) registerEndpoints(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Get("/version", svc.versionHandler)

	v1.Group("/get").
		Get("/transactions", svc.getTransactionsHandler)
}

func main() {
	svc := &service{}

	// Get CLI args.
	port := flag.Int("port", 8200, "Port number to host webserver")
	flag.Parse()

	app := fiber.New()
	svc.registerEndpoints(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", *port)))
}

// versionHandler handles requests to /version.
func (svc *service) versionHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"version": version})
}

func (svc *service) getTransactionsHandler(c *fiber.Ctx) error {
	// Load transactions.
	file, err := os.Open("transactions.csv")
	if err != nil {
		return c.JSON(fiber.Map{"error": fmt.Sprintf("unable to open transactions file: %v", err)})
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
