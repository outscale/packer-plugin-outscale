#!/bin/bash
set -e  # Exit immediately if a command fails

# Ensure wget is installed
if ! command -v wget >/dev/null 2>&1; then
    echo "Installing wget..."
    sudo apt-get update -y
    sudo apt-get install -y wget
fi

# Ensure libfuse2 is installed (required for AppImage execution)
if dpkg-query -W -f='${Status}\n' libfuse2t64:amd64 | grep -q "installed"; then
    echo "libfuse2 is already installed."
else
    echo "libfuse2 is not installed. Installing now..."
    sudo apt-get update -y
    sudo apt-get install -y libfuse2
fi


# Download the oapi-cli AppImage
echo "Downloading oapi-cli..."
wget -q https://github.com/outscale/oapi-cli/releases/latest/download/oapi-cli-x86_64.AppImage

# Make the file executable
chmod +x oapi-cli-x86_64.AppImage

# Move it to a system-wide location
echo "Installing oapi-cli..."
sudo mv oapi-cli-x86_64.AppImage /usr/local/bin/oapi-cli

# Ensure it's accessible from /usr/bin as well
sudo ln -sf /usr/local/bin/oapi-cli /usr/bin/oapi-cli

# Verify installation
echo "Verifying installation..."
oapi-cli --help
