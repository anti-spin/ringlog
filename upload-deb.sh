#!/bin/sh

# read DEB_USER, DEB_TOKEN, and DEB_REGISTRY from .env file
set -a
. ./.env
set +a

if [ $# -eq 0 ]; then
    echo "Error: Please provide a file path as an argument."
    echo "Usage: $0 <file_path>"
    exit 1
fi

file_path="$1"

if [ ! -f "$file_path" ]; then
    echo "Error: File '$file_path' does not exist."
    exit 1
fi

curl --user $DEB_USER:$DEB_TOKEN \
     --upload-file $file_path \
     $DEB_REGISTRY

# Check the curl exit status
if [ $? -eq 0 ]; then
    echo "Uploaded to debian registry $DEB_REGISTRY"
else
    echo "Error: File upload failed."
fi
