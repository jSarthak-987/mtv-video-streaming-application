# Use an official Node.js image as the base image for building
FROM node:18-alpine AS builder

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json (or yarn.lock) into the container
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code into the container
COPY . .

# Build the Next.js application
RUN npm run build

# Use a smaller Node.js image for the final image
FROM node:18-alpine AS runner

# Set the working directory in the container
WORKDIR /app

# Copy the built files from the builder stage
COPY --from=builder /app/.next .next
COPY --from=builder /app/public public
COPY --from=builder /app/package.json ./

# Install only production dependencies
RUN npm install --production

# Expose the port on which the app will run
EXPOSE 3000

# Start the Next.js application
CMD ["npm", "start"]
