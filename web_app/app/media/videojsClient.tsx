import React, { useEffect, useRef } from "react";
import videojs, { VideoJsPlayer, VideoJsPlayerOptions } from "video.js";
import "video.js/dist/video-js.css"; // Import default Video.js styles


// Define the props for the VideoJS component
interface VideoJSProps {
  options: VideoJsPlayerOptions;             // Options for configuring the Video.js player
  onReady?: (player: VideoJsPlayer) => void; // Optional callback to be called when the player is ready
}


// VideoJS component responsible for rendering and managing a Video.js player
export const VideoJS: React.FC<VideoJSProps> = ({ options, onReady }) => {
  const placeholderRef = useRef<HTMLDivElement | null>(null); // Ref to the placeholder div where the player will be mounted
  const playerRef = useRef<VideoJsPlayer | null>(null);       // Ref to store the Video.js player instance

  useEffect(() => {
    // Check if the player hasn't been initialized yet and the placeholder is available
    if (!playerRef.current && placeholderRef.current) {
      const placeholderEl = placeholderRef.current; // Reference to the placeholder element

      // Create a new <video-js> element and set its class for fullscreen
      const videoElement = document.createElement("video-js");
      videoElement.className = 'vjs-fullscreen'; // Add fullscreen class to ensure the player covers the entire screen

      // Append the video element to the placeholder div
      placeholderEl.appendChild(videoElement);

      // Initialize the Video.js player with the created video element and provided options
      const player = (playerRef.current = videojs(videoElement, options, () => {
        console.log("Player is ready"); // Log message indicating player readiness
        if (onReady) { // If an onReady callback is provided, call it with the player instance
          onReady(player);
        }
      }));

      // Set the player dimensions to cover the full window
      player.dimensions('100%', '100%');

      // Handler to adjust the player size when the window is resized
      const handleResize = () => {
        if (player) {
          player.width(window.innerWidth); // Update player width to match the window width
          player.height(window.innerHeight); // Update player height to match the window height
        }
      };

      // Add event listener for window resize to dynamically adjust player dimensions
      window.addEventListener('resize', handleResize);

      // Cleanup function to dispose of the player and remove the resize event listener when the component unmounts
      return () => {
        if (player) {
          player.dispose(); // Dispose of the Video.js player to free up resources
          playerRef.current = null; // Reset the player reference
        }
        window.removeEventListener('resize', handleResize); // Remove the resize event listener
      };
    }
  }, [options, onReady]); // Dependencies array to re-run the effect when options or onReady change

  // Render a div that acts as a placeholder for the Video.js player
  return (
    <div ref={placeholderRef} className="fixed inset-0"></div> // Set class for full screen coverage
  );
};

export default VideoJS;