package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rajputomsingh/relay/internal/config"
	"github.com/rajputomsingh/relay/internal/database"
	"github.com/rajputomsingh/relay/internal/handlers"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Relay stopped: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	startupCtx, cancelStartup := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	client, err := database.Connect(startupCtx, cfg.MongoURI)
	cancelStartup()

	if err != nil {
		return errors.New(
			"MongoDB initialization failed; check configuration, " +
				"Atlas network access, and database credentials",
		)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			log.Println("MongoDB disconnect failed")
		}
	}()

	eventsCollection := client.
		Database(cfg.MongoDatabase).
		Collection("events")

	indexCtx, cancelIndex := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	_, err = eventsCollection.Indexes().CreateOne(
		indexCtx,
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "idempotency_key", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetName("uniq_event_idempotency_key"),
		},
	)
	cancelIndex()

	if err != nil {
		return errors.New("failed to initialize event indexes")
	}

	mux := http.NewServeMux()

	healthHandler := handlers.NewHealthHandler(client)
	healthHandler.Register(mux)

	eventsHandler := handlers.NewEventsHandler(
		eventsCollection,
		cfg.RelayAPIKey,
	)
	eventsHandler.Register(mux)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	stopCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf(
			"Relay starting: environment=%s port=%s",
			cfg.Environment,
			cfg.Port,
		)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("HTTP server failed: %w", err)

	case <-stopCtx.Done():
		log.Println("Shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return errors.New("HTTP server shutdown timed out")
		}

		log.Println("HTTP server shut down gracefully")
		return nil
	}
}
