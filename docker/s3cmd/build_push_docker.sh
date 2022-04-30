#!/bin/bash
docker build -t zanderhavgaard/s3cmd:latest -f Dockerfile-s3cmd .
docker push zanderhavgaard/s3cmd:latest
