#!/bin/bash

# StudyForge Run Script

echo "Starting StudyForge..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go from https://golang.org/dl/"
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Error: .env file not found. Please copy .env.example to .env and configure it."
    exit 1
fi

# Create required directories
mkdir -p data uploads

# Load environment variables
set -a
source .env
set +a

# Check if HuggingFace API key is set
if [ -z "$HUGGINGFACE_API_KEY" ] || [ "$HUGGINGFACE_API_KEY" = "your_api_key_here" ]; then
    echo "Warning: HUGGINGFACE_API_KEY is not set in .env file"
    echo "Please add your Hugging Face API key to the .env file"
    exit 1
fi

# Run the server
echo "Starting server on http://$SERVER_HOST:$SERVER_PORT"
go run cmd/server/main.go
