package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jdotw/go-utils/log"
	"github.com/jdotw/go-utils/tracing"
	"github.com/jdotw/stock/pkg/transaction"
	itemapp "github.com/jdotw/stockonhand/internal/app/item"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func main() {
	serviceName := "stockonhand"

	// Logging and Tracing
	logger, metricsFactory := log.Init(serviceName)
	tracer := tracing.Init(serviceName, metricsFactory, logger)

	// Kafka Client
	seeds := []string{"localhost:9092"}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup("stock_on_hand"),
		kgo.ConsumeTopics("warehouse.stock.transaction.line_item.created"),
	)
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	// HTTP Router
	r := mux.NewRouter()

	// Item Service
	var itemService *itemapp.Service
	{
		repo, err := itemapp.NewGormRepository(context.Background(), os.Getenv("POSTGRES_DSN"), logger, tracer)
		if err != nil {
			logger.Bg().Fatal("Failed to create itemapp repository", zap.Error(err))
		}
		service := itemapp.NewService(repo, logger, tracer)
		endPoints := itemapp.NewEndpointSet(service, logger, tracer)
		itemapp.AddHTTPRoutes(r, endPoints, logger, tracer)
		itemService = &service
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

	go func() error {
		for {
			ctx := context.Background()
			fetches := cl.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				// All errors are retried internally when fetching, but non-retriable errors are
				// returned from polls so that users can notice and take action.
				logger.Bg().Fatal("Kafka fetch failed")
			}
			fetches.EachPartition(func(p kgo.FetchTopicPartition) {
				for _, record := range p.Records {
					var e transaction.TransactionLineItemCreatedEvent
					json.Unmarshal(record.Value, &e)
					err := (*itemService).UpdateStockOnHand(ctx, e.LineItem.ItemID, e.LineItem.Quantity)
					if err != nil {
						logger.Bg().Error("Failed to update stock on hand", zap.String("ID", e.LineItem.ID), zap.Error(err))
					}
					logger.Bg().Info("Received 'created' event for transaction line item", zap.String("ID", e.LineItem.ID))
				}
			})
		}
	}()

	// Select Loop
	select {}
}
