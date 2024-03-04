#!/bin/bash

# Define the name of the custom image
IMAGE_NAME="scratch-go-router:latest"
# Check if the image exists
if [[ "$(docker images -q $IMAGE_NAME 2> /dev/null)" == "" ]]; then
  echo "The router image '$IMAGE_NAME' does not exist. Building it..."
  
  # Build the image using the provided Dockerfile
  docker build --no-cache -t $IMAGE_NAME .

  if [ $? -ne 0 ]; then
    echo "Error: Failed to build the router image."
    exit 1
  fi
fi