#!/bin/bash

# Variables
IMAGE_NAME="forum-app"
CONTAINER_NAME="forum-container"
PORT=8081

# Stop and remove any existing container
docker stop $CONTAINER_NAME 2>/dev/null && docker rm $CONTAINER_NAME 2>/dev/null

# Prune only unused images and stopped containers (safer approach)
docker image prune -f
docker container prune -f

# Build the Docker image
docker build -t $IMAGE_NAME . && \
echo "Docker image built successfully." || \
{ echo "Failed to build Docker image."; exit 1; }

# Run the Docker container
docker run -d --name $CONTAINER_NAME -p $PORT:8081 $IMAGE_NAME && \
echo "Docker container is running on port $PORT." || \
{ echo "Failed to run Docker container."; exit 1; }