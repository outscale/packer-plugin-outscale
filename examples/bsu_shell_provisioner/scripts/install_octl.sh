#!/bin/bash
set -e  # Exit immediately if a command fails

# Ensure wget is installed
if ! command -v wget >/dev/null 2>&1; then
    echo "Installing wget..."
    sudo apt-get update -y
    sudo apt-get install -y wget
fi

# Download the octl AppImage
echo "Downloading octl..."
wget -o octl -q https://github.com/outscale/octl/releases/latest/download/octl_Linux_x86_64

# Make the file executable
chmod +x octl

# Move it to a system-wide location
echo "Installing octl..."
sudo mv octl /usr/local/bin

# Ensure it's accessible from /usr/bin as well
sudo ln -sf /usr/local/bin/octl /usr/bin/octl

# Verify installation
echo "Verifying installation..."
octl help
