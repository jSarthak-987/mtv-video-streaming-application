'use client';

import Heading from './heading';
import VideoList from './videoList';
import React from 'react';

// Main component for the Transcode page
export default function Transcode() {
    return (
      <React.Fragment>
        {/* Render the Heading component with a title prop */}
        <Heading title={"Video Transcoding App"} />
        
        {/* Render the VideoList component */}
        <VideoList />
      </React.Fragment>
    );
}