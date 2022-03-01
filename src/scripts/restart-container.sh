#!/bin/bash

# TODO: Check that the variables are actually supplied
if docker ps -q --filter "name=$1" | grep -q .
then
    # The docker container is currently running
    docker kill minitwit
    docker rm minitwit
fi

docker pull $2
docker run -d -p 8081:8081 --name $1 $2