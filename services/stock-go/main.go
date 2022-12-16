package main

import (
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	r, err := NewGormRepository(os.Getenv("POSTGRES_DSN"))
	ok(err)

	s, err := NewService(r)
	ok(err)

	//t := NewHTTPTransport(s) // 33k tps
	//t := NewGinTransport(s) // 33k tps
	t := NewFiberTransport(s) // 33k tps

	t.Serve()
}
