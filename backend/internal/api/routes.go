package api

import (
	"net/http"

	service "manhattan_tech_ventures/internal/services"

	"github.com/tus/tusd/v2/pkg/handler"
	"go.mongodb.org/mongo-driver/mongo"
)

// enableCORS is a middleware function that adds Cross-Origin Resource Sharing (CORS) headers
// to the HTTP response. This allows the server to handle requests from different origins,
// making it accessible from client applications hosted on other domains.
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers to allow all origins, methods, and specific headers.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests used by browsers to check CORS policy.
		if r.Method == http.MethodOptions {
			// Respond with a 204 No Content status to preflight requests.
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Pass the request to the next handler in the chain if not an OPTIONS request.
		next.ServeHTTP(w, r)
	})
}

// SetupRouter configures the HTTP router for the application by setting up routes
// for various endpoints such as file uploads, HLS streaming, and status updates.
// It integrates the TUS handler for file uploads and serves HLS media from GridFS.
func SetupRouter(tusHandler *handler.Handler, db *mongo.Database) *http.ServeMux {
	// Create a new ServeMux to handle routing.
	api := http.NewServeMux()

	// Set up an endpoint for server-sent events to stream status updates to clients.
	api.HandleFunc("/status/stream", func(w http.ResponseWriter, r *http.Request) {
		service.StatusStreamHandler(w, r)
	})

	// Set up TUS file upload endpoints, allowing clients to upload files to "/files/".
	// The uploaded files are managed by the TUS handler passed to the function.
	api.Handle("/files/", http.StripPrefix("/files/", tusHandler))
	api.Handle("/files", http.StripPrefix("/files", tusHandler))

	// Set up endpoints for serving HLS playlists (.m3u8) and segments (.ts) from GridFS.
	// CORS is enabled on these endpoints to allow requests from different origins.
	api.Handle("/hls", enableCORS(ServeM3U8(db)))    // Serve .m3u8 playlists
	api.Handle("/output/", enableCORS(ServeHLS(db))) // Serve HLS .ts segments

	// Serve static files from the "./web/static" directory for the root path.
	// This can be used for serving frontend assets like HTML, CSS, and JavaScript.
	api.Handle("/", http.FileServer(http.Dir("./web/static")))

	// Return the configured router.
	return api
}
