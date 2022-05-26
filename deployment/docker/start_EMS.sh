#!/bin/sh

myIPAddress=`ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $2}'`

docker run -e LOCAL_APIGW_SERVERS="http://${myIPAddress}:8001" -p 7001:80 --name website 192.168.1.12:5010/website:v1.0.8 > logs/website.log &

# apigw
docker run -p 8001:80 -e LOCAL_EXEC_SERVERS="http://${myIPAddress}:9001" --name apigw 192.168.1.12:5010/apigw:v1.1.1 > logs/apigw.log &

# exec
docker run -p 9001:80 -e DOWNLOAD_PORT="9001" -e MY_IP="${myIPAddress}" --name exec 192.168.1.12:5010/exec-python:v1.1.11 > logs/exec.log &

echo "Please access this URL: http://${myIPAddress}:7001/user/top"
