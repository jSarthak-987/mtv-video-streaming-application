'use client';

import React, { useState, useCallback, Suspense } from 'react';
import { useSearchParams } from 'next/navigation';
import VideoPlayer from './videojsClient';  // Import the VideoPlayer component
import videojs from 'video.js';             // Import Video.js for handling the video player


// Main component for rendering a video page with a video player
const VideoContent = () => {
    const searchParams = useSearchParams();        // Hook to access URL search parameters (query strings)
    const streamId = searchParams.get('streamId'); // Extract the 'streamId' parameter from the URL
    
    // State to store options for the Video.js player
    const [videoJsOptions, setVideoJsOptions] = useState<videojs.PlayerOptions | null>(null); 

    // Callback function that is triggered when the Video.js player is ready
    const handlePlayerReady = useCallback((player: videojs.Player) => {
        console.log('Player is ready:', player);    // Log a message when the player is ready
        player.on('play', () => {
            console.log('Video is playing');        // Log a message when the video starts playing
        });
    }, []);

    // Effect to set up the Video.js player options when the 'streamId' changes
    React.useEffect(() => {
        if (streamId) {
            // Set the options for the Video.js player including controls, fluid layout, and multiple source qualities
            setVideoJsOptions({
                controls: true, // Show player controls (play, pause, etc.)
                fluid: true, // Make the player responsive to window size changes
                sources: [
                    {
                        // URL for 480p stream
                        src: `http://localhost:8080/hls?quality=480p&stream_id=${streamId}`, 
                        type: 'application/x-mpegURL', // MIME type for HLS (HTTP Live Streaming)
                    },
                    {
                        // URL for 720p stream
                        src: `http://localhost:8080/hls?quality=720p&stream_id=${streamId}`,
                        type: 'application/x-mpegURL', // MIME type for HLS
                    },
                ],
            });
        }
    }, [streamId]); // Dependency array with 'streamId' to re-run effect when 'streamId' changes


    // Render the video player with options, or a loading message if options are not yet set
    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-black">
            {videoJsOptions ? (
                // Render the VideoPlayer component with the configured options and ready handler
                <VideoPlayer options={videoJsOptions} onReady={handlePlayerReady} />
            ) : (
                // Display a loading message while the video player is being set up
                <p className="text-white">Loading player...</p>
            )}
        </div>
    );
};


// Main component wrapped in Suspense
const VideoPage = () => (
    <Suspense fallback={<p className="text-white">Loading...</p>}>
        <VideoContent />
    </Suspense>
);

export default VideoPage;
