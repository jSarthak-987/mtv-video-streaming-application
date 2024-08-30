package service

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"manhattan_tech_ventures/internal/config"
	"manhattan_tech_ventures/internal/storage"

	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	"github.com/tus/tusd/v2/pkg/handler"
	"go.mongodb.org/mongo-driver/mongo"
)

// UploadStatus represents the status of an upload, including the filename and its current status.
// This struct is used to track the progress of uploads and send updates to clients.
type UploadStatus struct {
	Filename string // Name of the uploaded file
	Status   string // Current status of the file (e.g., "Uploading", "Uploaded")
}

var (
	// uploadStatusList keeps track of all upload statuses for reporting and monitoring purposes.
	uploadStatusList []UploadStatus
)

// HandleUpload initializes and handles the TUS upload process, including setting up the storage,
// creating the handler, and managing the completion of uploads. It also integrates a worker pool
// to process completed uploads.
func HandleUpload(storageService storage.Storage, dbClient *mongo.Database) *handler.Handler {

	// Retrieve the base path for uploads from the local storage service
	uploadDir := storageService.(*storage.LocalStorage).GetBasePath()

	// Load configuration settings
	conf := config.LoadConfig()

	// Create channels for jobs and results to manage the worker pool
	jobs := make(chan Job, 100)      // Buffered channel for jobs to process
	results := make(chan error, 100) // Buffered channel for results from workers

	// Start the worker pool to handle transcoding and uploading tasks
	go WorkerPool(jobs, results)

	// Set up file storage and locking mechanisms for TUS
	store := filestore.New(uploadDir)   // Use filestore for TUS storage
	locker := filelocker.New(uploadDir) // Use file locker to manage concurrent access

	// Compose the TUS store and locker into a store composer
	composer := handler.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	// Create a TUS handler with a configuration that includes the store composer
	// and enables notification on completed uploads.
	tusHandler, err := handler.NewHandler(handler.Config{
		BasePath:              "/files/", // Base path for handling file uploads
		StoreComposer:         composer,  // Composer that includes file storage and locking
		NotifyCompleteUploads: true,      // Enable notifications when uploads are complete
	})

	if err != nil {
		log.Fatalf("unable to create handler: %s", err) // Log and terminate if the handler cannot be created
	}

	go func() {
		for event := range tusHandler.CompleteUploads {
			uploadID := event.Upload.ID         // Get the unique ID of the completed upload
			filename := filepath.Base(uploadID) // Extract the filename from the upload ID

			// Update the upload status list to reflect the file has been uploaded
			uploadStatusList = append(uploadStatusList, UploadStatus{
				Filename: filename,
				Status:   "Uploaded",
			})

			// Send a status update to the client indicating the file has been uploaded
			SendStatusUpdateToClient(currClientChan, fmt.Sprintf("UC-%s:OK", filename))

			// Send a job to the worker pool for transcoding and further processing
			jobs <- Job{
				DBBucketName:   "media",                 // The GridFS bucket name in MongoDB
				UploadPath:     conf.UploadPath,         // Path where the uploaded file is stored
				TranscodedPath: conf.TranscodedFilePath, // Path where the transcoded files will be stored
				DBClient:       dbClient,                // MongoDB client used for interacting with GridFS
				Filename:       filename,                // Name of the file to be processed
				ClientChan:     currClientChan,          // Client channel for sending status updates
			}
		}
	}()

	// Middleware to handle incoming TUS uploads
	tusHandler.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Extract the filename from the request query, defaulting to "unknown_filename" if not provided
		filename := r.URL.Query().Get("filename")
		if filename == "" {
			filename = "unknown_filename"
		}

		// Update the upload status list to reflect the file is currently uploading
		uploadStatusList = append(uploadStatusList, UploadStatus{
			Filename: filename,
			Status:   "Uploading",
		})

		// Serve the request through the TUS handler
		tusHandler.ServeHTTP(w, r)
	}))

	// Return the configured TUS handler
	return tusHandler
}
