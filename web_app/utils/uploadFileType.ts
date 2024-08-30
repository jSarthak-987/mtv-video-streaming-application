// Import the Status enum representing various states of video processing, 
// type definitions for video files and file sizes
import { Status, UploadedVideoFile, FileSize } from "./videoFile"; 

// Define the interface for the Zustand store that manages uploaded video files
export interface UploadedVideoFileStore {
    files: UploadedVideoFile[]; // Array to hold all uploaded video files

    // Function to create a new file entry in the store
    createNewFile: (file: UploadedVideoFile) => Promise<UploadedVideoFile>; 
    // Accepts an UploadedVideoFile object and returns a promise that resolves with the added file

    // Function to synchronize the generated custom hash with the actual tus ID after upload
    syncFileIdWithTusId: (generateCustomHash: string, tusFilId: string) => Promise<UploadedVideoFile[]>;
    // Accepts the custom hash and the tus ID, returning a promise that resolves with the updated list of files

    // Function to update the upload file size (progress) of a specific file
    updateFileSize: (fileName: string, newUploadFileSize: FileSize) => void; 
    // Accepts a file name and new file size, updating the corresponding file's size in the store

    // Function to update the status of a specific file by tus ID
    updateStatus: (tusFileId: string, updateString: Status) => UploadedVideoFile | undefined; 
    // Accepts the tus file ID and a new status, updating the file's status and returning the updated file or undefined if not found
}
