#!/bin/bash
set -e

if ! [ -x "$(command -v wget)" ]; then
    sudo apt-get -y update
    sudo apt-get -y install wget
fi

if ! [ -x "$(command -v fuser)" ]; then
    sudo apt-get -y update
    sudo apt-get -y install libfuse2t64
    fuser -v 
fi

wget https://github.com/outscale/oapi-cli/releases/download/v0.7.0/oapi-cli-x86_64.AppImage
chmod a+x ./oapi-cli-x86_64.AppImage
sudo mv oapi-cli-x86_64.AppImage /usr/local/bin/oapi-cli
sudo cp /usr/local/bin/oapi-cli  /usr/bin/oapi-cli
oapi-cli --help
