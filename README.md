# Manhattan Tech Ventures - Video Transcoding App

A Next.js/Golang application for uploading, transcoding, and streaming videos. This app allows users to upload video files, monitor the upload and transcoding progress, and stream the videos once transcoded.

This project is designed to handle file uploads, stream HLS media, and interact with a MongoDB database. The server is built using Go version `go1.20`, and the front end is built using Next.js version 18. Both can be easily run inside individual Docker containers for consistent and isolated environments.

This project uses Docker Compose to orchestrate all the individual components, so it's recommended that you run the application using Docker Compose (Instructions are given after the backend section). For running the frontend or backend individually, run them either using Docker itself or outside of the Docker environment. Instructions are given in each of the individual sections.


==========================================================================

## Front End Application Structure (NextJS)

```
web_app/
│
├── app/                          # Main application directory for Next.js with App Router
│   ├── media/                    # Media-related components and pages
│   │   ├── page.tsx              # Main page for media (e.g., video rendering)
│   │   └── videojsClient.tsx     # Video.js player component
│   │
│   ├── transcode/                # Transcoding-related components and pages
│   │   ├── heading.tsx           # Component for heading in the transcoding page
│   │   ├── page.tsx              # Main page for transcoding
│   │   └── videoList.tsx         # Component for listing videos
│   │
│   ├── layout.tsx                # Global layout for the app
│   ├── globals.css               # Global CSS styles
│   └── page.tsx                  # Main entry page for the app
│
├── services/                     # Contains service logic similar to backend services
│   └── contextService.tsx        # Context and SSE logic for real-time updates
│
├── store/                        # Zustand stores or other state management files
│   └── videoUploadStore.ts       # Store for managing video upload state
│
├── utils/                        # Utility functions and types
│   ├── uploadFileType.ts         # Type definitions for file uploads
│   └── videoFile.ts              # Type definitions for video files

├── public/                       # Static files accessible by the client
│   └── images.svg                # Images used in the project
│
├── .eslintrc.json                # ESLint configuration
├── .gitignore                    # Git ignore file
├── next-env.d.ts                 # Next.js TypeScript environment types
├── next.config.js                # Next.js configuration file
├── tsconfig.json                 # TypeScript configuration file
├── postcss.config.js             # PostCSS configuration file
├── tailwind.config.js            # Tailwind CSS configuration file
├── package.json                  # NPM package file
└── README.md                     # Project documentation
```

## Features

- **Resumable File Uploads**: Supports resumable uploads using the tus protocol, allowing large video files to be uploaded with the ability to resume if interrupted.
- **Real-Time Status Updates**: Uses Server-Sent Events (SSE) to provide real-time status updates on the progress of video uploads and transcoding.
- **Multi-Quality Video Streaming**: Streams videos in multiple qualities (480p, 720p) using HLS (HTTP Live Streaming), providing a smooth viewing experience across different network conditions.
- **Responsive UI**: Built with Next.js and styled using Tailwind CSS for a modern, responsive user interface that works across devices.
- **State Management with Zustand**: Utilizes Zustand for efficient state management, keeping track of video file statuses, upload progress, and other application states.

## How It Works

1. **Uploading Videos**: Users can select video files to upload. The app uses `tus-js-client` to handle resumable uploads, ensuring reliability for large files.
2. **Tracking Progress**: The app displays the progress of each video upload and transcoding operation, updating the status in real-time via SSE.
3. **Transcoding**: After uploading, videos are transcoded into multiple resolutions using a backend service (e.g., FFmpeg) to provide different quality streams.
4. **Streaming**: Once transcoded, videos can be streamed directly from the app using a custom Video.js player, with options to switch between available qualities.

## Getting Started

### Prerequisites

- **Node.js**: Ensure you have Node.js installed (version 14 or above recommended).
- **Yarn** or **npm**: A package manager for installing dependencies.

### Setup Instructions

1. **Clone the Repository**

   ```bash
   git clone https://github.com/jSarthak-987/mtv-video-streaming-application.git
   cd manhattan-tech-ventures
   cd web_app
    ```

2. **Install Dependencies**

    ```
    npm install
    ```

3. **Running the Development Server**

    Start the Next.js development server with npm:

    ```
    npm run dev
    ```
    The app will be available at `http://localhost:3000/transcode`.

4. **Running the Backend Server**:

    Ensure your backend server (e.g., Go server handling file uploads and transcoding) is running at the specified API URL. The backend should handle endpoints for:

    - *File Uploads*: `POST /files/`
    - *Status Stream*: `GET /status/stream`
    - *HLS Streaming*: `GET /hls`

### Project Structure

- **pages/:** Contains the main pages of the application, including video upload and streaming interfaces.
    
- **services/:** Contains service files like `contextService` for SSE connections.
    
- **store/:** State management using Zustand for handling video files and statuses.

- **public/:** Static assets like icons and images.

### Key Components

1. **VideoPlayer:** A custom component using Video.js for playing videos with multiple quality options.

2. **Heading:** Displays the application title and handles file input for uploading.

3. **VideoList:** Displays a list of uploaded videos with their status and progress indicators.

4. **SSEStatusProvider:** Manages real-time updates for video statuses using Server-Sent Events.


### Troubleshooting

1. **CORS Issues:** Ensure your backend server has proper CORS configurations to allow requests from the frontend domain.

2. **Network Errors:** Check your backend server logs if uploads or status updates fail; ensure endpoints are correctly implemented and reachable.



==========================================================================

## Back End Application Structure (Golang)

```
backend/
│
├── cmd/
│   └── server/
│       └── main.go              # Entry point for the application
│
├── internal/
│   ├── api/
│   │   ├── handlers.go          # HTTP handlers for HLS Media Streaming (e.g., Serve .m3u8 and .ts files) 
│   │   └── routes.go            # Define all API routes
│   │
│   ├── storage/
│   │   ├── storage.go           # Interface for storage operations
│   │   ├── local_storage.go     # Local filesystem storage implementation
│   │   └── s3_storage.go        # (Implement For Optimized Storing and CDN distribution) S3 storage implementation
│   │
│   ├── database/
│   │   └── connection.go        # Database Connection Functions, for MongoDB
│   │
│   ├── services/
│   │   ├── video_service.go     # Business logic for handling video upload, transcoding, etc.
│   │   └── notification.go      # Logic for notifying users about upload status using SSE Events
|   |   └── transcoder.go        # Logic for transcoding videos using FFmpeg
│   │
│   └── config/
│       └── config.go            # Configuration loading and environment variables
│
├── web (for testing)/
│   ├── static/                  # Static files (JS, CSS, etc.)
│   ├── templates/               # HTML templates if using server-side rendering
│   └── videojs/                 # Video.js player-related files
│
├── scripts/
│   └── migration.sql            # Database migration scripts (if applicable, for postgres)
│
├── go.mod                       # Go module file
├── go.sum                       # Go module dependencies checksums
└── README.md                    # Project documentation
```

## Project Overview

### Features

- **File Uploads**: Supports resumable uploads using TUS protocol, allowing clients to upload large files efficiently.
- **HLS Streaming**: Transcodes uploaded videos into HLS format, generating `.m3u8` playlists and `.ts` segments for adaptive streaming.
- **MongoDB Integration**: Stores and retrieves HLS media files from MongoDB GridFS, providing scalable storage for media content.
- **Real-time Status Updates**: Uses Server-Sent Events (SSE) to provide clients with real-time status updates on upload and transcoding processes.

### Backend Functionality

The backend is responsible for handling file uploads, managing transcoding operations, and serving media files for streaming. Below are key components of the backend:

1. **TUS Upload Handler**:
   - Handles resumable file uploads using the TUS protocol.
   - Saves uploaded files to local storage, making them available for further processing.

2. **Transcoding Service**:
   - Uses `FFmpeg` to transcode uploaded videos into HLS format at multiple resolutions (e.g., 480p, 720p).
   - Generates `.m3u8` playlist files and `.ts` segments, which are stored in MongoDB GridFS.

3. **MongoDB GridFS**:
   - Manages storage of transcoded media files using GridFS, a specification for storing and retrieving large files in MongoDB.
   - Serves media files to clients on demand, supporting adaptive streaming via HLS.

4. **Status Updates via SSE**:
   - Provides real-time updates on file upload, transcoding, and storage operations using Server-Sent Events.
   - Allows clients to monitor the progress of their uploads and transcoding jobs in real-time.

### Prerequisites

- **Docker**: Ensure Docker is installed and running on your system. Download Docker from [Docker's official website](https://www.docker.com/products/docker-desktop).
- **Go (Optional)**: Required for local development. Download from [Go's official website](https://golang.org/dl/).

## Getting Started

### 1. Clone the Repository

Clone the repository to your local machine:

    ```bash
    git clone https://github.com/jSarthak-987/mtv-video-streaming-application.git
    cd manhattan-tech-ventures
    cd web_app
    ```

### 2. Build the Docker Image

Use the following command to build the Docker image:

```
docker build -t my-go-server .
```

This command builds a Docker image named my-go-server using the instructions defined in the Dockerfile.

### 3. Run the Docker Container

Run the Docker container using the image you just built:

```
docker run -d -p 8080:8080 --name go-server my-go-server
```

- **d**: Runs the container in `detached mode`.
- **p 8080:8080**: Maps port `8080` on your host machine to port `8080` in the container.
- **name go-server**: Names the container `go-server`.


### 4. Access the Application

You can now access your Go server at `http://localhost:8080`. Use a web browser or a tool like `curl`:

```
curl http://localhost:8080
```

### 5. Environment Variables

If your application requires environment variables, you can pass them to the Docker container in two ways:

--  **Option 1: Pass Environment Variables Directly**:
    Run the container with environment variables passed using the `-e` flag:

    ```
    docker run -d -p 8080:8080 --name go-server -e SERVER_ADDRESS=":8080" \
    -e MONGO_URI="mongodb://localhost:27017" my-go-server
    ```

-- **Option 2: Use an Environment File**
    Create a `.env` file with your environment variables:

    ```
    SERVER_ADDRESS=:8080
    MONGO_URI=mongodb://localhost:27017
    ```

    Then, run the container with the `.env` file:

    ```
    docker run -d -p 8080:8080 --name go-server --env-file .env my-go-server
    ```

### 6. Viewing Running Containers

To check if your container is running, use:

```
docker ps
```

This will list all running containers, including `go-server`.

### 7. Stopping and Removing the Container


To stop the container:

```
docker stop go-server
```

To remove the container:

```
docker rm go-server
```

### 8. Removing the Docker Image

If you need to remove the Docker image:

```
docker rmi my-go-server
```

## Troubleshooting
1. **Port Conflicts:** Ensure that port 8080 is not being used by another application.

2. **Docker Permissions:** If you encounter permission issues, try running commands with `sudo` on Unix-based systems.

3. **MongoDB Connection:** Ensure your MongoDB instance is running and accessible from within the Docker container.

## Notes
This setup assumes your Go server is built for a *Linux* environment to ensure compatibility with Docker containers.

Adjustments may be needed if running in different environments or with specific network configurations.



==========================================================================

## Running the Application Using Docker Compose (Recommended, especially for linux host)

This guide will help you set up and run both the Next.js frontend and Golang backend applications using Docker Compose for streamlined development and deployment.

### Prerequisites

- **Docker:** Ensure Docker is installed and running on your system. You can download Docker from Docker's official website.

- **Docker Compose:** Docker Compose is included with Docker Desktop, but if you're using Linux, you may need to install it separately.

### Setting Up the Project

1. **Clone the Repository**

First, clone the repository to your local machine:

```
git clone https://github.com/jSarthak-987/mtv-video-streaming-application.git
cd manhattan-tech-ventures
```

2. **Directory Structure**

Ensure your project structure looks like this:

```
manhattan-tech-ventures/
├── web_app/                # Next.js frontend project
├── backend/                # Golang backend project
├── docker-compose.yml      # Docker Compose file to orchestrate both services
├── README.md               # Instruction Manual, Project Description, etc.
└── .gitignore              # For ignoring files during git push
```

### Docker Compose Configuration

Your `docker-compose.yml` file is already set up to run both applications in their respective containers. Here is an overview of what it does:


- **Frontend (Next.js):**

    1. Builds the Next.js frontend from the `web_app` directory.
    2. Exposes port `3000` for accessing the frontend.
    3. Connects to the backend via the environment variable `NEXT_PUBLIC_API_URL` set to `http://backend:8080`.

- **Backend (Golang):**

    1. Builds the Go backend from the `backend` directory.
    2. Exposes port `8080` for API access and HLS streaming.
    3. Connects to a MongoDB service running in another container.

- **MongoDB:**

    1. A MongoDB service running on port `27017`, used for storing media files via GridFS.


### Running the Application

To start the application, follow these steps:

1. **Navigate to the Root Directory**

    Ensure you are in the root directory (`manhattan-tech-ventures`) where the `docker-compose.yml` file is located.

2. **Start the Services**

    Use Docker Compose to build and start all services:

    ```
    docker-compose up --build
    ```

    This command will:

    1. Build the Docker images for both the frontend and backend services.
    2. Start the Next.js frontend on `http://localhost:3000`.
    3. Start the Go backend on `http://localhost:8080`.
    4. Start MongoDB on `localhost:27017` within the Docker network.
    
3. **Access the Application**

    - ***Frontend:*** Visit `http://localhost:3000` in your browser to access the Next.js application.
    - ***Backend API:*** You can interact with the backend API at `http://localhost:8080.
   
4. **Stopping the Application**

    To stop the running containers, use:

    ```
    docker-compose down
    ```

    This command will stop and remove all containers started by Docker Compose.


## Troubleshooting

1. **Port Conflicts:** If ports `3000`, `8080`, or `27017` are already in use, you'll need to stop the services using these ports or adjust the port mappings in the `docker-compose.yml` file.

2. **Connection Issues:** Ensure your Docker network is correctly set up. If the frontend cannot connect to the backend, verify the `NEXT_PUBLIC_API_URL` environment variable in the frontend service configuration.

3. **Logs and Debugging:** Use `docker-compose logs -f` to view the logs from all services and troubleshoot issues as they occur.


## Environment Variables

To modify environment variables or configurations, update the respective sections in the `docker-compose.yml` file:

- **Frontend Environment:** Set in the frontend service under environment.
- **Backend Environment:** Set in the backend service under environment.

For example, you can set the backend service URL by modifying the `NEXT_PUBLIC_API_URL` environment variable in the frontend service section.


## Additional Notes
This setup assumes that all dependencies and configurations are correctly defined in the respective Dockerfile and docker-compose.yml.


________________________________________________________________________________

*Developed by Sarthak Joshi as part of Manhattan Tech Ventures Assignment.*
