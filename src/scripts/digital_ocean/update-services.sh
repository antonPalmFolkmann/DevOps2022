#!/bin/bash
docker pull antonfolkmann/minitwit:latest
docker service update --image antonfolkmann/minitwit:latest minitwit_web