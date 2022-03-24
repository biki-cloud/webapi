#!/bin/zsh



# ----------- apigw
echo -e "\n\n\n\n\n"
echo "################### apigw #####################"

echo "kubectl delete -f apigw.yml"
kubectl delete -f apigw.yml


# cd to project root
cd ../..


echo "rm -fr deployment/k8s/apigw-kube/apigw/apigw"
rm -fr deployment/k8s/apigw-kube/apigw/apigw

echo "GOOS=linux GOARCH=amd64 go build -o deployment/k8s/apigw-kube/apigw/apigw microservices/apigw/cmd/application/main/main.go"
GOOS=linux GOARCH=amd64 go build -o deployment/k8s/apigw-kube/apigw/apigw microservices/apigw/cmd/application/main/main.go

cd deployment/k8s/apigw-kube

#echo "docker rmi -f 192.168.1.12:5010/apigw:v1"
#docker rmi -f 192.168.1.12:5010/apigw:v1

# build from Dockerfile
echo "docker build -t 192.168.1.12:5010/apigw:v1 ."
docker build -t 192.168.1.12:5010/apigw:v1 .

echo "docker push 192.168.1.12:5010/apigw:v1 ."
docker push 192.168.1.12:5010/apigw:v1

# cd k8s
cd ..

echo "kubectl apply -f apigw.yml"
kubectl apply -f apigw.yml
