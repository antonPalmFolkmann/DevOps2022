#!/bin/bash

docker build -f docker/webserver.Dockerfile -t antonfolkmann/minitwit .
docker build -f docker/db.Dockerfile -t antonfolkmann/minitwitdb .