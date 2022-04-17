#!/bin/zsh


echo -e "\n\n\n\n\n"
# ----------- exec1
echo "######################### exec-python2 ###########################"
# pythonコンテナを元にしたexec pod
# pythonプログラムを持つexec pod

#echo "kubectl delete -f exec-python2.yml"
#kubectl delete -f exec-python1-1.yml

# cd to project root
cd ../..

echo "rm -fr deployment/k8s/exec-python/exec/exec"
rm -fr deployment/k8s/exec-python/exec/exec

echo "GOOS=linux GOARCH=amd64 go build -o deployment/k8s/exec-python/exec/exec microservices/exec/cmd/application/main/main.go"
GOOS=linux GOARCH=amd64 go build -o deployment/k8s/exec-python/exec/exec microservices/exec/cmd/application/main/main.go

cd deployment/k8s/exec-python

echo "docker rmi 192.168.1.12:5010/exec-python:v1.0.1"
docker rmi 192.168.1.12:5010/exec-python:v1.0.1

# build from Dockerfile
echo "docker build --no-cache -t 192.168.1.12:5010/exec-python:v1.0.4"
docker build --no-cache -t 192.168.1.12:5010/exec-python:v1.0.8 .

echo "docker push 192.168.1.12:5010/exec-python:v1.0.4"
docker push 192.168.1.12:5010/exec-python:v1.0.8

# cd k8s
cd ..

echo "kubectl delete -f exec-python2.yml"
kubectl delete -f exec-python2.yml

echo "kubectl apply -f exec-python2.yml"
kubectl apply -f exec-python2.yml
