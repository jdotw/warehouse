package main

import (
	"context"
	_ "embed"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jdotw/go-utils/log"
	"github.com/jdotw/go-utils/tracing"
	categoryapp "github.com/jdotw/stock/internal/app/category"
	itemapp "github.com/jdotw/stock/internal/app/item"
	locationapp "github.com/jdotw/stock/internal/app/location"
	transactionapp "github.com/jdotw/stock/internal/app/transaction"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func main() {
	serviceName := "stock"

	// Logging and Tracing
	logger, metricsFactory := log.Init(serviceName)
	tracer := tracing.Init(serviceName, metricsFactory, logger)

	// HTTP Router
	r := mux.NewRouter()

	// Create Kafka Client
	seeds := []string{"localhost:9092"}
	k, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
	)
	if err != nil {
		panic(err)
	}
	defer k.Close()

	// Category Service
	{
		repo, err := categoryapp.NewGormRepository(context.Background(), os.Getenv("POSTGRES_DSN"), logger, tracer)
		if err != nil {
			logger.Bg().Fatal("Failed to create categoryapp repository", zap.Error(err))
		}
		service := categoryapp.NewService(repo, logger, tracer)
		endPoints := categoryapp.NewEndpointSet(service, logger, tracer)
		categoryapp.AddHTTPRoutes(r, endPoints, logger, tracer)
	}

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

	// Location Service
	{
		repo, err := locationapp.NewGormRepository(context.Background(), os.Getenv("POSTGRES_DSN"), logger, tracer)
		if err != nil {
			logger.Bg().Fatal("Failed to create locationapp repository", zap.Error(err))
		}
		service := locationapp.NewService(repo, logger, tracer)
		endPoints := locationapp.NewEndpointSet(service, logger, tracer)
		locationapp.AddHTTPRoutes(r, endPoints, logger, tracer)
	}

	// Transaction Service
	{
		repo, err := transactionapp.NewGormRepository(context.Background(), os.Getenv("POSTGRES_DSN"), logger, tracer)
		if err != nil {
			logger.Bg().Fatal("Failed to create transactionapp repository", zap.Error(err))
		}
		service := transactionapp.NewService(repo, k, logger, tracer)
		endPoints := transactionapp.NewEndpointSet(service, logger, tracer)
		transactionapp.AddHTTPRoutes(r, endPoints, logger, tracer)
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
			httpPort = "8080"
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
