package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage provides a simple implementation of a storage system
// that uses the local filesystem to store, retrieve, and delete files.
type LocalStorage struct {
	BasePath string // Base path where files will be stored
}

// NewLocalStorage creates a new instance of LocalStorage with the specified base path.
// It initializes the storage with the provided path, which is used for all file operations.
func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{
		BasePath: basePath,
	}
}

// GetBasePath returns the base directory path used by the storage.
// This is useful for accessing or displaying the current storage location.
func (s *LocalStorage) GetBasePath() string {
	return s.BasePath
}

// Save writes the provided data to a file with the given filename within the storage's base path.
// It creates the base directory if it does not exist, and then creates the file and writes the data into it.
func (s *LocalStorage) Save(filename string, data io.Reader) (string, error) {
	// Ensure the base directory exists by creating it if necessary.
	err := os.MkdirAll(s.BasePath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create base directory: %v", err) // Return an error if directory creation fails.
	}

	// Construct the full file path within the base directory.
	filePath := filepath.Join(s.BasePath, filename)

	// Create the file at the specified path.
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err) // Return an error if file creation fails.
	}
	defer out.Close() // Ensure the file is closed after writing.

	// Write the data to the file from the provided reader.
	_, err = io.Copy(out, data)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err) // Return an error if writing to the file fails.
	}

	// Return the full path of the saved file.
	return filePath, nil
}

// Retrieve opens and returns a reader for the specified file within the storage's base path.
// This allows for reading the file's contents without loading the entire file into memory.
func (s *LocalStorage) Retrieve(filename string) (io.Reader, error) {
	// Construct the full file path within the base directory.
	filePath := filepath.Join(s.BasePath, filename)

	// Open the file for reading.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err) // Return an error if the file cannot be opened.
	}

	// Return the file as an io.Reader.
	return file, nil
}

// Delete removes the specified file from the storage's base path.
// This operation is permanent and will remove the file from the filesystem.
func (s *LocalStorage) Delete(filename string) error {
	// Construct the full file path within the base directory.
	filePath := filepath.Join(s.BasePath, filename)

	// Remove the file from the filesystem.
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err) // Return an error if the file cannot be deleted.
	}

	// Return nil if the file is successfully deleted.
	return nil
}
