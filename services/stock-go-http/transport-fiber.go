package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

type FiberTransport struct {
	r Repository
}

func NewFiberTransport(r Repository) FiberTransport {
	return FiberTransport{
		r: r,
	}
}

func (t FiberTransport) getCategories(c *fiber.Ctx) error {
	v, err := t.r.GetCategories(context.Background())
	if err != nil {
		panic(err)
	}
	return c.JSON(v)
}

func (t FiberTransport) Serve() {
	app := fiber.New()
	app.Get("/categories", t.getCategories)
	log.Fatal(app.Listen(":8080"))
}
