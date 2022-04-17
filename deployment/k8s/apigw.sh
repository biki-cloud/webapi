#!/bin/zsh



# ----------- apigw
echo -e "\n\n\n\n\n"
echo "################### apigw #####################"


# cd to project root
cd ../..


echo "rm -fr deployment/k8s/apigw-kube/apigw/apigw"
rm -fr deployment/k8s/apigw-kube/apigw/apigw

echo "GOOS=linux GOARCH=amd64 go build -o deployment/k8s/apigw-kube/apigw/apigw microservices/apigw/cmd/application/main/main.go"
GOOS=linux GOARCH=amd64 go build -o deployment/k8s/apigw-kube/apigw/apigw microservices/apigw/cmd/application/main/main.go

cd deployment/k8s/apigw-kube

echo "docker rmi 192.168.1.12:5010/apigw:v1.0.5"
docker rmi 192.168.1.12:5010/apigw:v1.0.5

# build from Dockerfile
echo "docker build --no-cache -t 192.168.1.12:5010/apigw:v1.0.6 ."
docker build --no-cache -t 192.168.1.12:5010/apigw:v1.0.8 .

echo "docker push 192.168.1.12:5010/apigw:v1.0.6 ."
docker push 192.168.1.12:5010/apigw:v1.0.8

# cd k8s
cd ..

echo "kubectl delete -f apigw.yml"
kubectl delete -f apigw.yml

echo "kubectl apply -f apigw.yml"
kubectl apply -f apigw.yml
