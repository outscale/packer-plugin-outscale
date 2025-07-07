#!/bin/bash
set -e

project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
EXAMPLES_DIR=$project_root/examples

make build
echo "example $EXAMPLES_DIR"
echo "root $project_root"
echo "projet $project_dir"

for f in $EXAMPLES_DIR/*
do
    if [ -d $f ]
    then
        cp packer-plugin-outscale $f/
        cd $f
        packer init .
        packer plugins install --path ./packer-plugin-outscale github.com/outscale/outscale
        packer fmt .
        packer validate .
        packer build .
        rm packer-plugin-outscale
        cd -
    fi
done

exit 0