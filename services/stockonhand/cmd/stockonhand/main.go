package main

import (
	"context"
	_ "embed"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jdotw/go-utils/log"
	"github.com/jdotw/go-utils/tracing"
	itemapp "github.com/jdotw/stockonhand/internal/app/item"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	serviceName := "stockonhand"

	// Logging and Tracing
	logger, metricsFactory := log.Init(serviceName)
	tracer := tracing.Init(serviceName, metricsFactory, logger)

	// HTTP Router
	r := mux.NewRouter()

	// Item Service
	{
		repo, err := itemapp.NewGormRepository(context.Background(), os.Getenv("POSTGRES_DSN"), logger, tracer)
		if err != nil {
			logger.Bg().Fatal("Failed to create itemapp repository", zap.Error(err))
		}
		service := itemapp.NewService(repo, logger, tracer)
		endPoints := itemapp.NewEndpointSet(service, logger, tracer)
		itemapp.AddHTTPRoutes(r, endPoints, logger, tracer)
	}

	// HTTP Mux
	m := tracing.NewServeMux(tracer)
	m.Handle("/metrics", promhttp.Handler()) // Prometheus
	m.Handle("/", r)

	// Start Transports
	go func() error {
		// HTTP
		httpHost := os.Getenv("HTTP_LISTEN_HOST")
		httpPort := os.Getenv("HTTP_LISTEN_PORT")
		if len(httpPort) == 0 {
			httpPort = "8083"
		}
		httpAddr := httpHost + ":" + httpPort
		logger.Bg().Info("Listening", zap.String("transport", "http"), zap.String("host", httpHost), zap.String("port", httpPort), zap.String("addr", httpAddr))
		err := http.ListenAndServe(httpAddr, m)
		logger.Bg().Fatal("Exit", zap.Error(err))
		return err
	}()

	// Select Loop
	select {}
}