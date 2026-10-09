package database

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(
		options.Client().
			ApplyURI(uri).
			SetAppName("relay").
			SetServerSelectionTimeout(5 * time.Second),
	)
	if err != nil {
		return nil, errors.New("failed to initialize MongoDB client")
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		cleanupCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		_ = client.Disconnect(cleanupCtx)
		return nil, errors.New("MongoDB is unreachable")
	}

	return client, nil
}
