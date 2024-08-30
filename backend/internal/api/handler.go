package api

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"manhattan_tech_ventures/internal/config"
	service "manhattan_tech_ventures/internal/services"

	"go.mongodb.org/mongo-driver/mongo"
)

// Load configuration settings using the LoadConfig function from the config package.
// This configuration will be used to determine paths for serving HLS streams and other settings.
var conf config.Config = config.LoadConfig()

// ServeM3U8 handles requests to serve .m3u8 files (HLS playlists) from MongoDB GridFS.
// It retrieves the desired quality and stream ID from query parameters, constructs the path to the .m3u8 file,
// and uses the ServeFileFromGridFS function to serve the file to the client.
func ServeM3U8(dbClient *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters from the URL
		queryParams := r.URL.Query()

		// Get the 'quality' parameter from the query, defaulting to "480p" if not provided
		quality := queryParams.Get("quality")
		if quality == "" {
			quality = "480p"
		}

		// Get the 'stream_id' parameter from the query
		streamId := queryParams.Get("stream_id")
		if streamId == "" {
			// If 'stream_id' is missing, respond with a bad request error
			http.Error(w, "Missing stream parameter", http.StatusBadRequest)
			return
		}

		// Construct the directory and file path for the .m3u8 playlist based on quality and stream ID
		m3u8Dir := filepath.Join(conf.TranscodedFilePath, quality, streamId, "m3u8")
		m3u8FilePath := "./" + filepath.Join(m3u8Dir, quality+".m3u8")

		// Serve the .m3u8 file from GridFS using the constructed file path
		service.ServeFileFromGridFS(w, r, dbClient, m3u8FilePath, "media")
	}
}

// ServeHLS handles requests to serve HLS segments (.ts files) from MongoDB GridFS.
// It parses the URL path to extract the quality, stream ID, and filename of the .ts segment,
// constructs the full path, and uses ServeFileFromGridFS to serve the file.
func ServeHLS(dbClient *mongo.Database) http.HandlerFunc {
	// Log the current client channel, primarily for debugging purposes
	log.Print(service.GetCurrClientChan())

	return func(w http.ResponseWriter, r *http.Request) {
		// Split the URL path to extract variables such as quality, stream ID, and filename
		parts := strings.Split(r.URL.Path, "/")

		// Check if the path is in the expected format; if not, return a bad request error
		if len(parts) < 5 {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		// Extract the quality, stream ID, and filename from the URL path
		quality := parts[2]
		streamID := parts[3]
		filename := parts[5]

		// Construct the full path to the .ts segment based on the extracted variables
		filePath := fmt.Sprintf("%s/%s/%s/ts/%s", conf.TranscodedFilePath, quality, streamID, filename)

		// Serve the .ts file from GridFS using the constructed file path
		service.ServeFileFromGridFS(w, r, dbClient, filePath, "media")
	}
}
