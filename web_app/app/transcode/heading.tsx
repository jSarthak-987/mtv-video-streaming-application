import { useRef } from 'react';
import * as tus from 'tus-js-client'; // Import tus client for resumable uploads
import { createHash } from 'crypto'; // Import createHash from Node.js crypto module for generating unique file hashes
import { useUploadedVideoFileStore } from '@store/videoUploadStore'; // Zustand store for managing video file state
import { Status } from '@utils/videoFile'; // Enum for file statuses

// Define the props for the Heading component
type HeaderProps = {
    title: string; // Title of the header, displayed at the top of the page
};

// Heading component for rendering the header with upload functionality
export default function Heading({ title }: HeaderProps) {
    const fileInputRef = useRef<HTMLInputElement>(null); // Reference to the hidden file input element
    const currDateTime: Date = new Date(); // Current date and time for generating unique hashes

    // Extract functions from the Zustand store for managing video file uploads
    const { updateVideoFileSize, createNewVideoFile, syncFileIdWithTusId, updateFileStatus } = useUploadedVideoFileStore((state) => ({
        updateVideoFileSize: state.updateFileSize, // Function to update file size progress
        createNewVideoFile: state.createNewFile, // Function to create a new file entry in the store
        syncFileIdWithTusId: state.syncFileIdWithTusId, // Function to sync custom IDs with Tus IDs
        updateFileStatus: state.updateStatus // Function to update the status of the file
    }));

    // Handler to simulate clicking the hidden file input when the upload button is clicked
    const handleButtonClick = () => {
        if (fileInputRef.current) {
            fileInputRef.current.click();
        }
    };

    // Handler for processing file input changes (file selection)
    const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
        if (event.target.files && event.target.files.length > 0) {

            const file = event.target.files[0]; // Get the first selected file
            // Generate a unique hash based on the file name and current time
            const dateTimeStr = currDateTime.getDate() + ':' + currDateTime.getTime();
            const hashString = file.name + '' + dateTimeStr;

            // Function to generate a SHA-256 hash from a string
            const generateCustomHash = ((input: string) => {
                return createHash('sha256').update(input).digest('hex');
            })(hashString);

            // Create a new video file entry in the store with initial data
            await createNewVideoFile({
                id: generateCustomHash, // Use the generated hash as a temporary ID
                videoTitle: file.name,
                uploadFileSize: {
                    totalSize: file.size,
                    uploadedSize: 0
                },
                status: Status.UPLOAD_STARTED, // Initial status set to 'UPLOAD_STARTED'
                progressValue: '0%' // Initial progress value
            });

            // Create a new tus upload instance with the selected file
            const upload: tus.Upload = new tus.Upload(file, {
                endpoint: 'http://localhost:8080/files/', // Tus server endpoint for file uploads
                onError: function (error: Error) {
                    // Error handling for upload failures
                    console.error("Failed because: " + error.message);
                    console.error("Error details: ", error.stack);
                    const fileId = upload.url?.split('//').pop()?.split('/').pop(); // Extract file ID from URL
                    if (fileId) {
                        updateFileStatus(fileId, Status.UPLOAD_FAILURE); // Update status to 'UPLOAD_FAILURE' if ID exists
                    }
                },
                onProgress: function (bytesUploaded, bytesTotal) {
                    // Update file size progress in the store during upload
                    updateVideoFileSize(generateCustomHash, { totalSize: bytesTotal, uploadedSize: bytesUploaded });
                },
                onSuccess: async function () {
                    const fileId = upload.url?.split('//').pop()?.split('/').pop(); // Extract the final file ID from the upload URL
                    if (fileId) {
                        // Sync the temporary ID with the final ID from the tus server
                        await syncFileIdWithTusId(generateCustomHash, fileId); // Update the ID in the store
                        updateFileStatus(fileId, Status.UPLOAD_SUCCESS); // Update status to 'UPLOAD_SUCCESS'
                    }
                    console.log('Download %s from %s', upload.file, upload.url); // Log the successful upload details
                }
            });
            upload.start(); // Start the tus upload
            event.target.value = ''; // Reset the file input value
        }
    };

    return (
        <div className="md:flex border-b-2 border-b-gray-900 pb-8 md:items-center md:justify-between ml-8 mr-8 mt-8">
            <div className="min-w-0 flex-1">
                {/* Display the title passed as a prop */}
                <h2 className="text-2xl font-bold leading-7 text-white sm:truncate sm:text-3xl sm:tracking-tight">
                    {title}
                </h2>
            </div>
            <div className="mt-4 flex md:ml-4 md:mt-0">
                {/* Button to trigger file input click */}
                <button
                    type="button"
                    onClick={handleButtonClick}
                    className="transition duration-300 ease-in-out ml-3 inline-flex items-center rounded-md bg-indigo-600 px-6 py-3 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500"
                >
                    Upload
                </button>

                {/* Hidden file input for selecting files */}
                <input
                    type="file"
                    ref={fileInputRef} // Ref to access the input programmatically
                    onChange={handleFileChange} // Handle file selection
                    style={{ display: 'none' }} // Hide the input element
                />
            </div>
        </div>
    );
}
