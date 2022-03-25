#!/bin/zsh

files/apigw1/apigw -port 8001 > log/apigw1.log &
files/apigw2/apigw -port 8002 > log/apigw2.log &
files/exec1/exec -port 9001 > log/exec1.log &
files/exec2/exec -port 9002 > log/exec2.log &
files/website1/website -port 7001 > log/website1.log &
