#!/bin/bash

# TODO:
# npm run build /path/to/package.json/ --output=./client/build

docker build -f docker/webserver.Dockerfile -t antonfolkmann/minitwit .
docker build -f docker/db.Dockerfile -t antonfolkmann/minitwitdb .
