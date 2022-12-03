package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	r, err := NewGormRepository(context.Background(), os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatalln(err)
	}

	// t := NewHTTPTransport(r)
	t := NewGinTransport(r)

	t.Serve()
}
