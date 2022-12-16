package transport

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

type fiberTransport struct {
	s Service
}

func NewFiberTransport(s Service) Transport {
	return fiberTransport{
		s: s,
	}
}

func (t fiberTransport) getCategories(c *fiber.Ctx) error {
	v, err := t.s.GetCategories(context.Background())
	ok(err)
	return c.JSON(v)
}

func (t fiberTransport) Serve() {
	app := fiber.New()
	app.Get("/categories", t.getCategories)
	log.Fatal(app.Listen(":8080"))
}
