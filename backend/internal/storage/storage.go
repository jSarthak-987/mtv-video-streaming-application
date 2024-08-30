package storage

import "io"

// Storage defines an interface for file storage operations, including saving and retrieving files.
// This interface provides a standard set of methods that any storage implementation must fulfill,
// allowing for flexibility in swapping out different storage backends (e.g., local, cloud, etc.)
type Storage interface {
	// Save writes the provided data to a file with the specified filename.
	// It returns the full path of the saved file or an error if the operation fails.
	// The data to be saved is provided as an io.Reader, allowing for flexible data sources.
	Save(filename string, file io.Reader) (string, error)

	// Retrieve opens and returns a reader for the specified file, allowing for streaming
	// the file's contents. This method returns an io.Reader to read the file's data
	// or an error if the file cannot be accessed.
	Retrieve(filename string) (io.Reader, error)
}
