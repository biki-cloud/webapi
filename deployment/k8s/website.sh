#!/bin/zsh




echo -e "\n\n\n\n\n"
# ----------- website
echo "################ website #####################"


# cd to project root
cd ../..

echo "rm -fr deployment/k8s/website-kube/website/website"
rm -fr deployment/k8s/website-kube/website/website

echo "GOOS=linux GOARCH=amd64 go build -o deployment/k8s/website-kube/website/website microservices/website/cmd/application/main/main.go"
GOOS=linux GOARCH=amd64 go build -o deployment/k8s/website-kube/website/website microservices/website/cmd/application/main/main.go

rm -fr deployment/k8s/website-kube/website/ui
cp -r microservices/website/cmd/application/ui deployment/k8s/website-kube/website/ui

cd deployment/k8s/website-kube

echo "docker rmi 192.168.1.12:5010/website:v1.0.5"
docker rmi 192.168.1.12:5010/website:v1.0.5

echo "docker build --no-cache -t 192.168.1.12:5010/website:v1.0.6 ."
docker build --no-cache -t 192.168.1.12:5010/website:v1.0.6 .

echo "docker push 192.168.1.12:5010/website:v1.0.6"
docker push 192.168.1.12:5010/website:v1.0.6

# cd k8s
cd ..

echo "kubectl apply -f website.yml"
kubectl apply -f website.yml
