#!/bin/sh

myIPAddress=`ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $2}' | sed -n 1p`
echo "Your IPAddress: ${myIPAddress}"

echo ""
websiteCommand="docker run -e LOCAL_APIGW_SERVERS=http://${myIPAddress}:8001 -p 7001:80 --name website bikibiki/website:v1.0.9"
echo "${websiteCommand}"
${websiteCommand} > logs/website.log &

echo ""
# apigw
apigwCommand="docker run -p 8001:80 -e LOCAL_EXEC_SERVERS=http://${myIPAddress}:9001 --name apigw bikibiki/apigw:v1.1.1"
echo "${apigwCommand}"
${apigwCommand} > logs/apigw.log &


echo ""
# exec
execCommand="docker run -p 9001:80 -e DOWNLOAD_PORT=9001 -e MY_IP=${myIPAddress} --name exec bikibiki/exec-python:v1.0.20"
echo "${execCommand}"
${execCommand} > logs/exec.log &


echo ""
echo ""
echo "Please access this URL: http://${myIPAddress}:7001/user/top"
