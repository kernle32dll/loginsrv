#!/bin/bash

docker build -t kernle32dll/loginsrv . ;
docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD" ;
docker push kernle32dll/loginsrv ;
