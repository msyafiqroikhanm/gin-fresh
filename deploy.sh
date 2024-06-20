#!/bin/bash

# Function to print messages with color
print_message() {
  local color_code="$1"
  local message="$2"
  echo -e "\e[${color_code}m${message}\e[0m"
}

# Change to the project directory
print_message "33" "========== Changing Directory =========="
cd /home/syafiq/jxb-eprocurement || { print_message "31" "Failed to change directory"; exit 1; }
print_message "32" "Directory changed to /home/syafiq/jxb-eprocurement"

# Load .env file
print_message "33" "========== Loading Environment Variables =========="
if [ -f .env ]; then
  print_message "32" "Loading environment variables from .env file"
  export $(grep -v '^#' .env | xargs)
else
  print_message "31" ".env file not found"
  exit 1
fi

# Determine the branch based on the ENV variable
print_message "33" "========== Determining Branch =========="
if [ "$ENV" == "production" ]; then
  BRANCH="main"
  print_message "32" "Environment is production, pulling from branch: $BRANCH"
elif [ "$ENV" == "development" ]; then
  BRANCH="develop"
  print_message "32" "Environment is development, pulling from branch: $BRANCH"
else
  print_message "31" "Environment variable ENV is not set or has an invalid value"
  exit 1
fi

# Pull the latest changes from the git repository
print_message "33" "========== Pulling Latest Changes =========="
git pull origin $BRANCH || { print_message "31" "Git pull failed"; exit 1; }
print_message "32" "Pulled the latest changes from branch: $BRANCH"

# Remove the existing main executable
print_message "33" "========== Removing Existing Executable =========="
rm -f main || { print_message "31" "Failed to remove main executable"; exit 1; }
print_message "32" "Removed the existing main executable"

# Rebuild main.go
print_message "33" "========== Rebuilding main.go =========="
go build -o main main.go || { print_message "31" "Build failed"; exit 1; }
print_message "32" "Rebuilt main.go successfully"

# Restart the systemd service
print_message "33" "========== Restarting Service =========="
sudo systemctl restart go-jxb-eprocurement.service || { print_message "31" "Failed to restart service"; exit 1; }
print_message "32" "Service restarted successfully"

print_message "33" "========== Deployment Completed =========="
print_message "32" "Deployment completed successfully."
