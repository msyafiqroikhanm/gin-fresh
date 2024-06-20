#!/bin/bash

# Load .env file
if [ -f .env ]; then
    echo "Loading environment variables from .env file"
    export $(grep -v '^#' .env | xargs)
else
    echo ".env file not found"
    exit 1
fi

# Change to the project directory
echo "Changing directory to /home/syafiq/jxb-eprocurement"
cd /home/syafiq/jxb-eprocurement || { echo "Failed to change directory"; exit 1; }

# Determine the branch based on the ENV variable
if [ "$ENV" == "production" ]; then
    BRANCH="main"
    echo "Environment is production, pulling from branch: $BRANCH"
elif [ "$ENV" == "development" ]; then
    BRANCH="develop"
    echo "Environment is development, pulling from branch: $BRANCH"
else
    echo "Environment variable ENV is not set or has an invalid value"
    exit 1
fi

# Pull the latest changes from the git repository
echo "Pulling the latest changes from the Git repository"
git pull origin $BRANCH || { echo "Git pull failed"; exit 1; }

# Remove the existing main executable
echo "Removing the existing main executable"
rm -f main || { echo "Failed to remove main executable"; exit 1; }

# Rebuild main.go
echo "Rebuilding main.go"
go build -o main main.go || { echo "Build failed"; exit 1; }

# Restart the systemd service
echo "Restarting the go-jxb-eprocurement.service"
sudo systemctl restart go-jxb-eprocurement.service || { echo "Failed to restart service"; exit 1; }

echo "Deployment completed successfully."
