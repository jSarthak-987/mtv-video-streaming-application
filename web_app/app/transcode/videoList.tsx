import React, { useEffect, useState, useMemo, useRef } from 'react';
import { ChevronRightIcon, PlayCircleIcon } from '@heroicons/react/20/solid';
import { ReceiveSSEStatus } from '@/services/contextService';         // Custom hook to receive SSE (Server-Sent Events) status updates
import { useUploadedVideoFileStore } from '@store/videoUploadStore';  // Zustand store for managing video uploads
import { UploadedVideoFile, Status } from '@utils/videoFile';         // Types for uploaded videos and statuses
import { useRouter } from 'next/navigation';


// CSS classes for different status colors
const statuses = {
  Gray: 'text-gray-700 bg-gray-100/10',
  Blue: 'text-sky-500 bg-sky-100/10',
  Green: 'text-green-400 bg-green-400/10',
  Red: 'text-rose-400 bg-rose-400/10',
};

// Status categories grouped by color
const GreenStatuses: Status[] = [
  Status.UPLOAD_STARTED,
  Status.UPLOAD_SUCCESS,
  Status.TRANSCODE_STARTED,
  Status.TRANSCODE_480_SUCCESS,
  Status.TRANSCODE_720_SUCCESS,
];

const RedStatuses: Status[] = [
  Status.TRANSCODE_FAILURE,
  Status.UPLOAD_FAILURE,
];

const BlueStatus: Status[] = [
  Status.TRANSCODE_SUCCESS,
];

// Utility function to concatenate class names
function classNames(...classes: string[]) {
  return classes.filter(Boolean).join(' ');
}

// Main component for displaying the list of uploaded videos
export default function VideoList() {
  const router = useRouter();
  const uploadSuccess = useRef(false);          // Ref to track whether an upload was successful
  const sseStatusMessage = ReceiveSSEStatus();  // Hook to receive real-time updates from SSE
  const { files, updateStatus } = useUploadedVideoFileStore((state) => ({
    files: state.files,                         // List of uploaded files from Zustand store
    updateStatus: state.updateStatus,           // Function to update file status in store
  }));
  
  // State to manage the list of video files
  const [videoFileListItem, setVideoFileListItem] = useState<UploadedVideoFile[]>([]);


  // useMemo to calculate the updated list of video files with progress values
  const updatedItems: UploadedVideoFile[] = useMemo(() => {
    return files.map((fileItem) => {

      // Calculate upload progress as a percentage
      const percentageValue = ((fileItem.uploadFileSize.uploadedSize * 100) / fileItem.uploadFileSize.totalSize).toFixed(2);
      const progressValue = String(percentageValue) + '%';

      return {
        id: fileItem.id,
        videoTitle: fileItem.videoTitle,
        status: fileItem.status,
        progressValue,
        uploadFileSize: {
          totalSize: fileItem.uploadFileSize.totalSize,
          uploadedSize: fileItem.uploadFileSize.uploadedSize,
        },
      };
    });
  }, [files]);


  // Effect to update the list of video files when the updated items change
  useEffect(() => {
    // Only update state if there are changes
    if (JSON.stringify(videoFileListItem) !== JSON.stringify(updatedItems)) {
      setVideoFileListItem(updatedItems);
    }
  }, [updatedItems, videoFileListItem]);

  
  // Effect to handle incoming SSE status updates
  useEffect(() => {
    if (sseStatusMessage && sseStatusMessage.fileId && sseStatusMessage.message && sseStatusMessage.statusCategory) {
      let status: Status;

      // Determine status based on SSE message and category
      if (sseStatusMessage.message !== 'OK') {
        status = sseStatusMessage.statusCategory === 'UC' ? Status.UPLOAD_FAILURE : Status.TRANSCODE_FAILURE;
      } else {
        // Map status category to specific statuses
        switch (sseStatusMessage.statusCategory) {
          case 'UC':
            uploadSuccess.current = true;
            status = Status.UPLOAD_SUCCESS;
            break;
          case 'TS':
            uploadSuccess.current = true;
            status = Status.TRANSCODE_STARTED;
            break;
          case 'T4':
            uploadSuccess.current = true;
            status = Status.TRANSCODE_480_SUCCESS;
            break;
          case 'T7':
            uploadSuccess.current = true;
            status = Status.TRANSCODE_720_SUCCESS;
            break;
          case 'TC':
            uploadSuccess.current = true;
            status = Status.TRANSCODE_SUCCESS;
            break;
          default:
            status = Status.UPLOAD_STARTED;
            break;
        }
      }

      // Update the status of the video file in the Zustand store
      updateStatus(sseStatusMessage.fileId, status);
    }
    console.log(sseStatusMessage ? sseStatusMessage : "Empty"); // Debugging log for SSE messages
  }, [sseStatusMessage, updateStatus]);

  // Function to navigate to the video page
  const openVideoPage = (streamId: string) => {
    router.push(`/media?streamId=${streamId}`); // Navigate to the media page with the stream ID
  };

  // Render the list of video files with status indicators and controls
  return (
    <div>
      <ul role="list" className="divide-y divide-white/5 ml-8 mr-8 mt-4">
        {videoFileListItem.map((deployment) => (
          <li key={deployment.id} className={classNames(BlueStatus.includes(deployment.status) ? 'bg-blue-950 hover:bg-blue-900' : 'bg-black', "transition ease-in-out hover:bg-gray-950 hover:cursor-pointer relative flex items-center space-x-4 py-2 px-2 my-1 overflow-hidden rounded-lg shadow-slate-900")}>
            <div className="min-w-0 flex-auto px-4 py-2 sm:p-2">
              <div className="flex items-center gap-x-3">
                <span className={classNames(
                  GreenStatuses.includes(deployment.status) ? statuses.Green :
                    RedStatuses.includes(deployment.status) ? statuses.Red :
                      BlueStatus.includes(deployment.status) ? statuses.Blue : '',
                  'relative flex h-2 w-2'
                )}>
                  {GreenStatuses.includes(deployment.status) && (
                    <span className="animate-ping rounded-full bg-current absolute inline-flex h-full w-full opacity-75" />
                  )}
                  <span className="h-full w-full inline-flex rounded-full bg-current" />
                </span>
                <h2 className="min-w-0 text-l font-bold leading-2 text-gray-300">
                  <div className="flex gap-x-2">
                    <span className="truncate">{deployment.videoTitle}</span>
                  </div>
                </h2>
                {/* Play icon to navigate to the video player page when clicked */}
                <PlayCircleIcon onClick={() => openVideoPage(deployment.id)} aria-hidden="true" className={classNames(deployment.status === Status.TRANSCODE_SUCCESS ? 'opacity-100 visible': 'opacity-0 collapse', "transition-opacity h-5 w-5 flex-none text-gray-400")} />
              </div>
              <div className="mt-3 flex items-center gap-x-2.5 text-xs leading-5 text-gray-400">
                <p className="whitespace-nowrap">
                  {deployment.status === Status.UPLOAD_STARTED ? 'Uploading' :
                    deployment.status === Status.UPLOAD_SUCCESS ? 'Uploaded' :
                      deployment.status === Status.UPLOAD_FAILURE ? 'Upload Failure' : 'Uploaded'}
                </p>
                {uploadSuccess.current && (
                  <div className="flex items-center">
                    <ChevronRightIcon aria-hidden="true" className="h-5 w-5 flex-none text-gray-400" />
                    <p className="whitespace-nowrap">
                      {deployment.status === Status.TRANSCODE_STARTED ? 'Transcoding' :
                        deployment.status === Status.TRANSCODE_SUCCESS ? 'Transcoded' :
                          deployment.status === Status.TRANSCODE_FAILURE ? 'Transcode Failure' : 'Transcoding'}
                    </p>
                  </div>
                )}
              </div>
            </div>
            <div className='flex flex-row pr-4'>
              {/* Progress bar for upload status */}
              <div aria-hidden="true" className={classNames('transition-width duration-150 ease-in-out', [Status.UPLOAD_STARTED, Status.UPLOAD_FAILURE].includes(deployment.status) ? 'w-24' : 'w-5')}>
                <div className="overflow-hidden rounded-full bg-gray-900">
                  <div style={{ width: deployment.progressValue }} className={classNames(deployment.status === Status.UPLOAD_FAILURE ? 'bg-red-600' : deployment.status === Status.UPLOAD_STARTED ? 'bg-indigo-600' : 'bg-green-600', "translation-width duration-75 ease-in-out h-2")} />
                </div>
              </div>

              {/* Progress bar for transcoding status */}
              <div aria-hidden="true" className={classNames('transition-width duration-150 ease-in-out', [Status.UPLOAD_SUCCESS, Status.TRANSCODE_STARTED, Status.TRANSCODE_FAILURE, Status.TRANSCODE_480_SUCCESS, Status.TRANSCODE_720_SUCCESS].includes(deployment.status) ? 'w-24' : 'w-5', "ml-4")}>
                <div className="overflow-hidden rounded-full bg-gray-900">
                  <div style={{
                    width: (
                      deployment.status === Status.TRANSCODE_480_SUCCESS ? '50%' : // Set width based on transcoding progress
                        deployment.status === Status.TRANSCODE_720_SUCCESS ? '95%' : 
                          deployment.status === Status.TRANSCODE_SUCCESS ? '100%' : 
                            '0%') // Default width if no progress
                  }} className={classNames(deployment.status === Status.TRANSCODE_FAILURE ? 'bg-red-600' : deployment.status >= Status.TRANSCODE_STARTED && deployment.status < Status.TRANSCODE_SUCCESS ? 'bg-indigo-600' : 'bg-green-600', "translation-width duration-75 ease-in-out h-2")} />
                </div>
              </div>
            </div>
          </li>
        ))}
      </ul>

      {/* 
      // Conditional rendering for VideoPlayer dialog component, if needed in future
      <Dialog open={open} onClose={() => setOpen(false)} className="relative z-10">
        <DialogBackdrop
          transition
          className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity data-[closed]:opacity-0 data-[enter]:duration-300 data-[leave]:duration-200 data-[enter]:ease-out data-[leave]:ease-in"
        />

        <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <DialogPanel
              transition
              className="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all data-[closed]:translate-y-4 data-[closed]:opacity-0 data-[enter]:duration-300 data-[leave]:duration-200 data-[enter]:ease-out data-[leave]:ease-in sm:my-8 sm:w-full sm:max-w-sm sm:p-6 data-[closed]:sm:translate-y-0 data-[closed]:sm:scale-95"
            >
              {videoJsOptions ? (
                <VideoPlayer options={videoJsOptions} onReady={handlePlayerReady} />
              ) : (
                <p>Loading player...</p>
              )}
            </DialogPanel>
          </div>
        </div>
      </Dialog>
      */}
    </div>
  );
}