#!/bn/bash

docker build -t devopsitu/minitwitimage -f Dockerfile-minitwit .
docker build -t devopsitu/itusqlimage -f Dockerfile-mysql .

docker push devopsitu/minitwitimage:latest
docker push devopsitu/itusqlimage:latest
