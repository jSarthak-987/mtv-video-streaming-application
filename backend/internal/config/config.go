package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds configuration values for the application.
// These values are loaded from environment variables, providing flexibility for different environments.
type Config struct {
	ServerAddress      string // Address where the server will listen, e.g., ":8080"
	MongoURI           string // URI for connecting to MongoDB, e.g., "mongodb://localhost:27017"
	UploadPath         string // Path where uploaded files will be stored
	TranscodedFilePath string // Path where transcoded files will be stored
	WorkerProcessCount string // Number of worker processes for handling jobs concurrently
	DBName             string // Name of the MongoDB database used for storing media files
}

// LoadConfig loads configuration values from environment variables or uses default values if not set.
// It first attempts to load variables from a .env file if it exists, providing a convenient way
// to set environment variables for development.
func LoadConfig() Config {
	// Load .env file if available
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found") // Log if the .env file is not found; continue using system environment variables
	}

	// Return a Config struct populated with values from environment variables or default values
	return Config{
		ServerAddress:      getEnv("SERVER_ADDRESS", ":8080"),                // Default server address
		MongoURI:           getEnv("MONGO_URI", "mongodb://localhost:27017"), // Default MongoDB URI
		UploadPath:         getEnv("UPLOAD_PATH", "./uploads"),               // Default upload path
		TranscodedFilePath: getEnv("TRANSCODE_PATH", "./output"),             // Default transcoded files path
		WorkerProcessCount: getEnv("WP_COUNT", "2"),                          // Default number of worker processes
		DBName:             getEnv("DB_NAME", "hls_media"),                   // Default MongoDB database name
	}
}

// getEnv retrieves the value of an environment variable given by key.
// If the environment variable is not set, it returns the specified default value.
func getEnv(key, defaultValue string) string {

	// Check if the environment variable exists and return its value if it does
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	// Return the default value if the environment variable is not set
	return defaultValue
}
