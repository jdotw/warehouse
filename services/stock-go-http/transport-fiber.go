package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

type fiberTransport struct {
	r Repository
}

func NewFiberTransport(r Repository) Transport {
	return fiberTransport{
		r: r,
	}
}

func (t fiberTransport) getCategories(c *fiber.Ctx) error {
	v, err := t.r.GetCategories(context.Background())
	if err != nil {
		panic(err)
	}
	return c.JSON(v)
}

func (t fiberTransport) Serve() {
	app := fiber.New()
	app.Get("/categories", t.getCategories)
	log.Fatal(app.Listen(":8080"))
}
