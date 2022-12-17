package transport

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jdotw/warehouse/services/stock-go/internal/service"
	"github.com/jdotw/warehouse/services/stock-go/internal/util"
)

type fiberTransport struct {
	s service.Service
}

func NewFiberTransport(s service.Service) Transport {
	return fiberTransport{
		s: s,
	}
}

func (t fiberTransport) getCategories(c *fiber.Ctx) error {
	v, err := t.s.GetCategories(context.Background())
	util.Ok(err)
	return c.JSON(v)
}

func (t fiberTransport) Serve() {
	app := fiber.New()
	app.Get("/categories", t.getCategories)
	log.Fatal(app.Listen(":8080"))
}
