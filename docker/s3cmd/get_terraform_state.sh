#!/bin/bash

command=$1

if [ $command = "ls" ]; then
    docker run \
        --rm \
        -it \
        -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
        -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
        zanderhavgaard/s3cmd \
        s3cmd ls --recursive "s3://$SPACE_NAME"
elif [ $command = "get" ]; then
    docker run \
        --rm \
        -it \
        -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
        -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
        --mount type=bind,source="$(pwd)",target=/home/arch/shared \
        zanderhavgaard/s3cmd \
        s3cmd get --force "s3://$SPACE_NAME/$STATE_FILE"
elif [ $command = "it" ]; then
    docker run \
        --rm \
        -it \
        -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
        -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
        --mount type=bind,source="$(pwd)",target=/home/arch/shared \
        zanderhavgaard/s3cmd
else
    echo "Please use a valid command: ls, get, it"
fi
