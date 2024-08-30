package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongoDB establishes a connection to a MongoDB instance using the provided URI.
// It returns the MongoDB client, a context for further database operations, and an error if the connection fails.
// The connection is established with a timeout of 10 seconds to prevent hanging in case of network issues.
func ConnectMongoDB(uri string) (*mongo.Client, context.Context, error) {

	// Create a context with a timeout to limit the duration of the connection attempt.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is canceled to release resources if the connection attempt is done.

	// Attempt to connect to MongoDB using the provided URI and options.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {

		// Return an error if the connection to MongoDB fails.
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Return the MongoDB client, a background context for database operations, and no error.
	return client, context.Background(), nil
}
