#!/bin/bash
containers=$(docker ps -a -q -f "name=minitwit")
if [ -v $containers ]
then 
    echo "No containers to remove"
else
    docker rm $containers
fi