package main

import (
	"os"

	"github.com/jdotw/warehouse/services/stock-go/internal/repository"
	"github.com/jdotw/warehouse/services/stock-go/internal/service"
	"github.com/jdotw/warehouse/services/stock-go/internal/transport"
	"github.com/jdotw/warehouse/services/stock-go/internal/util"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	r, err := repository.NewGormRepository(os.Getenv("POSTGRES_DSN"))
	util.Ok(err)

	s, err := service.NewService(r)
	util.Ok(err)

	//t := transport.NewHTTPTransport(s) // 33k tps
	//t := transport.NewGinTransport(s) // 33k tps
	t := transport.NewFiberTransport(s) // 33k tps

	t.Serve()
}
