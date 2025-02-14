#!/bin/bash

# Variables
IMAGE_NAME="forum-app"
CONTAINER_NAME="forum-container"
PORT=8080

# Stop and remove any existing container
docker stop $CONTAINER_NAME 2>/dev/null && docker rm $CONTAINER_NAME 2>/dev/null

# Build the Docker image
docker build -t $IMAGE_NAME . && \
echo "Docker image built successfully." || \
{ echo "Failed to build Docker image."; exit 1; }

# Run the Docker container
docker run -d --name $CONTAINER_NAME -p $PORT:8080 $IMAGE_NAME && \
echo "Docker container is running on port $PORT." || \
{ echo "Failed to run Docker container."; exit 1; }