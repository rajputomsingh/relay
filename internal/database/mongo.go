package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	if uri == "" {
		return nil, fmt.Errorf("MongoDB URI is empty; check MONGODB_URI in .env")
	}

	client, err := mongo.Connect(
		options.Client().
			ApplyURI(uri).
			SetAppName("relay").
			SetServerSelectionTimeout(5 * time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MongoDB client: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		cleanupCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		_ = client.Disconnect(cleanupCtx)

		return nil, fmt.Errorf("MongoDB ping failed: %w", err)
	}

	return client, nil
}
