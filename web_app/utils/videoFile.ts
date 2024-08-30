// Define a type for file size, including total and uploaded sizes
export type FileSize = {
    totalSize: number,      // The total size of the file in bytes
    uploadedSize: number    // The size of the file that has been uploaded in bytes
};

// Define an enumeration for the various statuses a file can have during upload and transcoding
export enum Status {
    'UPLOAD_STARTED',          // Status when the upload has started
    'UPLOAD_SUCCESS',          // Status when the upload has successfully completed
    'UPLOAD_FAILURE',          // Status when the upload has failed
    'TRANSCODE_STARTED',       // Status when transcoding has started
    'TRANSCODE_480_SUCCESS',   // Status when 480p transcoding has successfully completed
    'TRANSCODE_720_SUCCESS',   // Status when 720p transcoding has successfully completed
    'TRANSCODE_SUCCESS',       // Status when transcoding has fully completed
    'TRANSCODE_FAILURE'        // Status when transcoding has failed
}

// Define a type for uploaded video files, which includes details about the file and its status
export type UploadedVideoFile = {
    id: string,                   // Optional ID for the video file, typically used for identifying the file
    videoTitle: string;            // The title of the video file
    uploadFileSize: FileSize,      // An object containing information about the total and uploaded sizes
    status: Status,                // The current status of the video file, based on the Status enum
    // uploadTime: string,         // (Optional) Time when the upload was initiated or completed
    // environment: 'Preview' | 'Production', // (Optional) Environment where the video is being processed, could be 'Preview' or 'Production'
    progressValue: string          // A string representing the progress of the upload/transcoding as a percentage
};
