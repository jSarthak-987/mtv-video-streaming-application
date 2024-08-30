package main

import (
	"log"
	"manhattan_tech_ventures/internal/api"
	"manhattan_tech_ventures/internal/config"
	database "manhattan_tech_ventures/internal/database"
	services "manhattan_tech_ventures/internal/services"
	"manhattan_tech_ventures/internal/storage"
	"net/http"
)

// main is the entry point of the application. It initializes the necessary services,
// sets up the router, and starts the HTTP server.
func main() {
	// Load configuration settings from environment variables or default values.
	cfg := config.LoadConfig()

	// Connect to MongoDB using the provided URI from the configuration.
	client, ctx, dberr := database.ConnectMongoDB(cfg.MongoURI)
	if dberr != nil {
		log.Fatalf("Error connecting to MongoDB: %v", dberr) // Log and terminate if the database connection fails.
	}

	// Access the specified MongoDB database using the database name from the configuration.
	db := client.Database(cfg.DBName)

	// Ensure that the MongoDB client disconnects properly when the application terminates.
	defer client.Disconnect(ctx)

	// Initialize local storage for file uploads using the base path from the configuration.
	storageService := &storage.LocalStorage{BasePath: cfg.UploadPath}

	// Set up the TUS upload handler using the storage service and MongoDB client.
	// This handler manages file uploads and integrates with the transcoding and storage services.
	tusHandler := services.HandleUpload(storageService, db)

	// Set up the HTTP router with the configured routes and TUS handler.
	// The router manages endpoints for uploads, HLS streaming, and status updates.
	http.Handle("/", api.SetupRouter(tusHandler, db))

	// Log that the server is running.
	log.Default().Printf("Server Running")

	// Start the HTTP server on the specified address from the configuration.
	err := http.ListenAndServe(cfg.ServerAddress, nil)
	if err != nil {
		log.Fatalf("unable to listen: %s", err) // Log and terminate if the server cannot start.
	} else {
		log.Default().Printf("Server Running") // Log that the server is running (redundant, already logged above).
	}
}
