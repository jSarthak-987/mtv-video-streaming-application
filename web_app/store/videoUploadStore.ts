import { create } from 'zustand';
import { UploadedVideoFileStore } from '@utils/uploadFileType'; // Import the type definition for the store
import { Status } from '@utils/videoFile'; // Import the Status enum for file status management


// Zustand store for managing the state of uploaded video files
export const useUploadedVideoFileStore = create<UploadedVideoFileStore>((set, get) => ({
  files: [], // Initialize the store with an empty array of files

  // Function to synchronize the generated custom hash with the tus ID after upload
  syncFileIdWithTusId: (generateCustomHash, tusFilId) => 
    new Promise((resolve, reject) => {
      try {
        // Update the file ID from the custom hash to the actual tus ID
        set((state) => ({
          files: state.files.map(fileObj => fileObj.id === generateCustomHash 
            ? { ...fileObj, id: tusFilId }  // Replace the temporary ID with the tus ID
            : fileObj                       // Leave other files unchanged
          )
        }));
        resolve(get().files);               // Resolve the promise with the updated list of files
      } catch(e) {
        reject(e);                          // Reject the promise if an error occurs
      }
    }),


  // Function to create a new file entry in the store
  createNewFile: (file) =>
    new Promise((resolve, reject) => {
      try {
        set((state) => ({ files: [...state.files, file] })); // Add the new file to the list
        resolve(file);
      } catch (error) {
        reject(error);
      }
    }),

    
  // Function to update the upload file size (progress) of a specific video by ID
  updateFileSize: (videoId, newFileSize) =>
    set((state) => ({
      files: state.files.map(fileObj => fileObj.id === videoId
          ? { ...fileObj, uploadFileSize: newFileSize } // Update the file size if IDs match
          : fileObj // Leave other files unchanged
      )
    })),

  // Function to update the status of a specific video by tus ID
  updateStatus: (tusId, newStatus: Status) => {
    set((state) => ({
      files: state.files.map(fileObj => fileObj.id === tusId
          ? { ...fileObj, status: newStatus } // Update the status if IDs match
          : fileObj // Leave other files unchanged
      )
    }));
    // Return the updated file object matching the given tus ID
    return get().files.find(fileObj => fileObj.id === tusId);
  }
}));
