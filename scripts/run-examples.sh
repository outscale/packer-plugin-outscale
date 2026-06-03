#!/bin/bash
set -e

project_dir=$(cd "$(dirname "$0")" && pwd)
project_root=$(cd "$project_dir"/.. && pwd)
EXAMPLES_DIR=$project_root/examples
PLUGIN_BINARY=$project_root/packer-plugin-outscale

prepare() {
    example_dir=$1
    packer_target=$2
    var_file="$example_dir/variables.auto.pkrvars.hcl"

    cp "$PLUGIN_BINARY" "$example_dir/"

    cd "$example_dir"
    packer init "$packer_target"
    packer plugins install --path ./packer-plugin-outscale github.com/outscale/outscale
    packer fmt "$packer_target"
    packer validate -var-file="$var_file" "$packer_target"

    rm "$example_dir/packer-plugin-outscale"
}

build() {
    example_dir=$1

    cp "$PLUGIN_BINARY" "$example_dir/"

    cd "$example_dir"
    packer build .

    rm "$example_dir/packer-plugin-outscale"
}

make build
echo "examples $EXAMPLES_DIR"

# Init and validate examples
for f in "$EXAMPLES_DIR"/*
do
    if [ -d "$f" ]; then
        prepare "$f" .
    elif [ -f "$f" ] && [[ "$f" == *.pkr.hcl ]]; then
        prepare "$EXAMPLES_DIR" "$f"
    fi
done

# Build examples
for f in "$EXAMPLES_DIR"/*
do
    if [ -d "$f" ]; then
        build "$f"
    fi
done

exit 0
