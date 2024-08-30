'use client';

import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';


// Define the structure of a status message received from the SSE server
export type StatusMessage = {
  message: 'OK' | string,   // The status message, which can be 'OK' or other strings indicating errors or other states
  fileId: string,          // The ID of the file related to this status message
  statusCategory: 'UC' | 'TC' | 'T4' | 'T7' | 'TS' // Status categories representing different stages of upload and transcoding
}


// Define the context type that will hold the status message
type SSEStatusContextType = {
  status: StatusMessage | null; // The status message or null if not yet received
};

// Create the context with an undefined initial value
const SSEStatusContext = createContext<SSEStatusContextType | undefined>(undefined);

// Provider component that manages the SSE connection and provides the status to its children
export const SSEStatusProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [status, setStatus] = useState<StatusMessage | null>(null); // State to hold the current status message

  useEffect(() => {
    // Initialize a new EventSource connection to receive SSE updates
    const eventSource = new EventSource("http://localhost:8080/status/stream");

    // Event handler for receiving messages from the SSE connection
    eventSource.onmessage = function (event) {
      const newStatus = event.data; // Get the data from the event
      console.log(newStatus);

      // Parse the status message from the event data
      const statusMessage = newStatus?.split(':').pop(); // Extract the actual message part
      const statusType = newStatus?.split(':')[0].split('-')[0]; // Extract the status type (e.g., 'UC')
      const fileId = newStatus?.split(':')[0].split('-')[1]; // Extract the file ID associated with the message

      // Construct a status message object based on the parsed data
      const statusObj: StatusMessage = {
        message: statusMessage,
        fileId: fileId,
        statusCategory: statusType
      }

      setStatus(statusObj); // Update the status state with the new message
    };

    // Error handler for the SSE connection
    eventSource.onerror = function (event) {
      console.error("EventSource failed:", event); // Log the error
      eventSource.close(); // Close the SSE connection on error to prevent further issues
    };

    // Cleanup function to close the SSE connection when the component unmounts
    return () => {
      eventSource.close(); // Ensure the connection is closed when the component is removed
    };
  }, []);

  // Provide the current status value to any children components through context
  return (
    <SSEStatusContext.Provider value={{ status }}>
      {children}
    </SSEStatusContext.Provider>
  );
};

// Custom hook to access the current SSE status from the context
export const ReceiveSSEStatus = () => {
  const context = useContext(SSEStatusContext); // Access the SSEStatusContext
  if (context === undefined) {
    // If the hook is used outside of a SSEStatusProvider, throw an error
    throw new Error('useSSEStatus must be used within a SSEStatusProvider');
  }
  return context.status; // Return the current status from the context
};
