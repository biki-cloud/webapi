#!/bin/sh

docker stop website
docker rm website

docker stop apigw
docker rm apigw

docker stop exec
docker rm exec

echo "Cleaned EMS containers"
