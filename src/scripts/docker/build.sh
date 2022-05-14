#!/bin/bash

rm -rf client/build
mkdir -p client
cd src/frontend 
npm run build
mv build ../../client/build
cd ../..

docker build -f docker/webserver.Dockerfile -t antonfolkmann/minitwit .
docker build -f docker/db.Dockerfile -t antonfolkmann/minitwitdb .
