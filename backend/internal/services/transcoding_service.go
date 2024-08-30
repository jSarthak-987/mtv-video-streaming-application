package service

import (
	"bytes"
	"fmt"
	"log"
	"manhattan_tech_ventures/internal/config"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

// maxWorkersEnv retrieves the maximum number of worker processes from the configuration file.
// This value is used to control the number of concurrent worker goroutines in the WorkerPool.
var maxWorkersEnv = config.LoadConfig().WorkerProcessCount

// Job represents a unit of work for the worker pool. It contains all the necessary
// information to process a video file, including paths, database client, and the channel
// to communicate status updates to the client.
type Job struct {
	DBBucketName   string          // Name of the GridFS bucket in MongoDB
	UploadPath     string          // Path where the original uploaded files are stored
	TranscodedPath string          // Path where transcoded files will be stored
	Filename       string          // Name of the original video file
	DBClient       *mongo.Database // MongoDB client used for GridFS operations
	ClientChan     chan string     // Channel for sending status updates back to the client
}

// WorkerPool initializes a pool of worker goroutines that process jobs from the jobs channel.
// Each worker transcodes videos and uploads them to GridFS, sending results to the results channel.
func WorkerPool(jobs chan Job, results chan error) {
	var wg sync.WaitGroup

	// Convert the maxWorkersEnv string to an integer
	maxWorkers, err := strconv.Atoi(maxWorkersEnv)
	if err != nil {
		panic(err) // Panic if maxWorkersEnv is not a valid number
	}

	// Start the specified number of worker goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				// Send a status update indicating the start of transcoding
				SendStatusUpdateToClient(job.ClientChan, fmt.Sprintf("TS-%s:OK", job.Filename))

				// Perform video transcoding and handle potential errors
				err := TranscodeVideo(job.DBClient, job.DBBucketName, job.UploadPath, job.TranscodedPath, job.Filename, job.ClientChan)

				// Send status updates based on the success or failure of the transcoding
				if err != nil {
					SendStatusUpdateToClient(job.ClientChan, fmt.Sprintf("TF-%s:%v", job.Filename, err))
				} else {
					SendStatusUpdateToClient(job.ClientChan, fmt.Sprintf("TC-%s:OK", job.Filename))
				}
				results <- err // Send the result (error or nil) to the results channel
			}
		}()
	}

	wg.Wait()      // Wait for all workers to complete
	close(results) // Close the results channel once all workers are done
}

// TranscodeVideo performs the transcoding of a video into 480p and 720p HLS streams.
// It creates necessary directories, runs FFmpeg commands for transcoding, and uploads the
// resulting HLS files to GridFS. Status updates are sent back to the client through a channel.
func TranscodeVideo(dbClient *mongo.Database, gridFSBucketName string, filePath string, outputfilePath string, originalFilename string, clientChanParam chan string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2) // Channel to collect errors from transcoding goroutines

	// Define the input and output paths for transcoding
	inputFullPath := filepath.Join(filePath, originalFilename)
	m3u8Output480FullPath := filepath.Join(outputfilePath, "480p", originalFilename, "m3u8")
	m3u8output720FullPath := filepath.Join(outputfilePath, "720p", originalFilename, "m3u8")
	tsOutput480FullPath := filepath.Join(outputfilePath, "480p", originalFilename, "ts")
	tsOutput720FullPath := filepath.Join(outputfilePath, "720p", originalFilename, "ts")

	// Create necessary directories for output files
	m3u8File480pErr := os.MkdirAll(m3u8Output480FullPath, os.ModePerm)
	if m3u8File480pErr != nil {
		log.Fatalf("failed to create base directory: %v", m3u8File480pErr)
		return m3u8File480pErr
	}

	tsFile480pErr := os.MkdirAll(tsOutput480FullPath, os.ModePerm)
	if tsFile480pErr != nil {
		log.Fatalf("failed to create base directory: %v", tsFile480pErr)
		return tsFile480pErr
	}

	m3u8File720pErr := os.MkdirAll(m3u8output720FullPath, os.ModePerm)
	if m3u8File720pErr != nil {
		log.Fatalf("failed to create base directory: %v", m3u8File720pErr)
		return m3u8File720pErr
	}

	tsFile720pErr := os.MkdirAll(tsOutput720FullPath, os.ModePerm)
	if tsFile720pErr != nil {
		log.Fatalf("failed to create base directory: %v", tsFile720pErr)
		return tsFile720pErr
	}

	// Define paths for the m3u8 playlists and .ts segments
	m3u8Output480p := filepath.Join(m3u8Output480FullPath, "480p.m3u8")
	m3u8Output720p := filepath.Join(m3u8output720FullPath, "720p.m3u8")
	tsOutput480p := filepath.Join(tsOutput480FullPath, "480p_%03d.ts")
	tsOutput720p := filepath.Join(tsOutput720FullPath, "720p_%03d.ts")

	// Define FFmpeg commands for transcoding to 480p and 720p HLS streams
	cmd480p := exec.Command("ffmpeg", "-i", inputFullPath,
		"-vf", "scale=-2:480",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-hls_segment_filename", tsOutput480p,
		"-hls_base_url", tsOutput480FullPath+"\\",
		"-f", "hls",
		m3u8Output480p)

	cmd720p := exec.Command("ffmpeg", "-i", inputFullPath,
		"-vf", "scale=-2:720",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-hls_segment_filename", tsOutput720p,
		"-hls_base_url", tsOutput720FullPath+"\\",
		"-f", "hls",
		m3u8Output720p)

	wg.Add(2)

	// Run FFmpeg command for 480p in a separate goroutine
	go func() {
		defer wg.Done()
		var stderr480p bytes.Buffer
		cmd480p.Stderr = &stderr480p
		err480 := cmd480p.Run()
		if err480 != nil {
			log.Printf("FFmpeg 480p error: %s", stderr480p.String())
			errChan <- fmt.Errorf("failed to transcode 480p: %v", err480)
		} else {
			SendStatusUpdateToClient(clientChanParam, fmt.Sprintf("T4-%s:OK", originalFilename))
		}
	}()

	// Run FFmpeg command for 720p in a separate goroutine
	go func() {
		defer wg.Done()
		var stderr720p bytes.Buffer
		cmd720p.Stderr = &stderr720p
		err720 := cmd720p.Run()
		if err720 != nil {
			log.Printf("FFmpeg 720p error: %s", stderr720p.String())
			errChan <- fmt.Errorf("failed to transcode 720p: %v", err720)
		} else {
			SendStatusUpdateToClient(clientChanParam, fmt.Sprintf("T7-%s:OK", originalFilename))
		}
	}()

	wg.Wait()      // Wait for both transcoding processes to complete
	close(errChan) // Close the error channel after all goroutines are done

	// Collect files to upload to GridFS
	var filesToUpload []string

	fileReadErr := filepath.Walk(outputfilePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filesToUpload = append(filesToUpload, ".\\"+path) // Add each file to the list
		}
		return nil
	})

	if fileReadErr != nil {
		log.Fatalf("Error traversing directory: %v", fileReadErr)
	}

	// Upload each file to GridFS
	for _, filePath := range filesToUpload {
		err := UploadFileToGridFS(dbClient, filePath, gridFSBucketName)
		if err != nil {
			log.Printf("Error uploading file %s: %v", filePath, err)
		}
	}

	// Combine errors from both transcoding processes
	var combinedError error
	for err := range errChan {
		if combinedError == nil {
			combinedError = err
		} else {
			combinedError = fmt.Errorf("%v; %v", combinedError, err)
		}
	}

	return combinedError
}
