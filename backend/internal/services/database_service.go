package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NormalizePath replaces backslashes with forward slashes to ensure consistent path formatting.
// This function is particularly useful for normalizing file paths on different operating systems.
func NormalizePath(path string) string {
	return strings.ReplaceAll(path, `\`, `/`)
}

// UploadFileToGridFS uploads a file from the local filesystem to MongoDB GridFS.
// It takes the MongoDB database, the file path of the file to upload, and the GridFS bucket name as parameters.
// The function opens the file, creates an upload stream in the specified GridFS bucket, and copies the file's contents
// into the GridFS bucket. If any step fails, it returns an error.
func UploadFileToGridFS(db *mongo.Database, filePath string, bucketName string) error {
	// Open the local file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a new GridFS bucket with the specified name
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(bucketName))
	if err != nil {
		return fmt.Errorf("failed to create GridFS bucket: %v", err)
	}

	// Normalize the file path by trimming the output directory prefix
	relPath := NormalizePath(strings.TrimPrefix(filePath, "./output/"))

	// Open an upload stream for the file in the GridFS bucket
	uploadStream, err := bucket.OpenUploadStream(relPath)
	if err != nil {
		return fmt.Errorf("failed to open upload stream: %v", err)
	}
	defer uploadStream.Close()

	// Copy the file's contents to the GridFS upload stream
	if _, err := io.Copy(uploadStream, file); err != nil {
		return fmt.Errorf("failed to upload file to GridFS: %v", err)
	}

	return nil
}

// ServeFileFromGridFS serves a file stored in MongoDB GridFS to the client over HTTP.
// It takes the HTTP response writer, request, MongoDB database, the filename to serve, and the GridFS bucket name as parameters.
// The function retrieves the file from GridFS using the filename, sets the appropriate content type based on the file extension,
// and writes the file's contents to the response. If the file is not found or any error occurs, it sends an appropriate HTTP error response.
func ServeFileFromGridFS(w http.ResponseWriter, r *http.Request, db *mongo.Database, filename string, bucketName string) {
	// Normalize the filename to ensure consistent path formatting
	normalizedFilePath := NormalizePath(filename)

	// Create a new GridFS bucket with the specified name
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(bucketName))
	if err != nil {
		http.Error(w, "Failed to create GridFS bucket", http.StatusInternalServerError)
		log.Printf("Failed to create GridFS bucket: %v", err)
		return
	}

	// Open a download stream for the file in the GridFS bucket using its name
	downloadStream, err := bucket.OpenDownloadStreamByName(normalizedFilePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		log.Printf("File not found in GridFS: %s", normalizedFilePath)
		return
	}
	defer downloadStream.Close()

	// Set the appropriate Content-Type header based on the file extension
	if filepath.Ext(normalizedFilePath) == ".m3u8" {
		w.Header().Set("Content-Type", "application/text")
	} else if filepath.Ext(normalizedFilePath) == ".ts" {
		w.Header().Set("Content-Type", "video/vnd.dlna.mpeg-tts")
	}

	// Copy the contents of the file from GridFS to the HTTP response writer
	if _, err := io.Copy(w, downloadStream); err != nil {
		http.Error(w, "Failed to serve file", http.StatusInternalServerError)
		log.Printf("Failed to serve file from GridFS: %v", err)
	}
}
