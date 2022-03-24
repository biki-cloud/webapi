#!/bin/zsh




echo -e "\n\n\n\n\n"
# ----------- website
echo "################ website #####################"

echo "kubectl delete -f website.yml"
kubectl delete -f website.yml


# cd to project root
cd ../..

echo "rm -fr deployment/k8s/website-kube/website/website"
rm -fr deployment/k8s/website-kube/website/website

echo "GOOS=linux GOARCH=amd64 go build -o deployment/k8s/website-kube/website/website microservices/website/cmd/application/main/main.go"
GOOS=linux GOARCH=amd64 go build -o deployment/k8s/website-kube/website/website microservices/website/cmd/application/main/main.go

rm -fr deployment/k8s/website-kube/website/ui
cp -r microservices/website/cmd/application/ui deployment/k8s/website-kube/website/ui

cd deployment/k8s/website-kube


echo "docker build -t 192.168.1.12:5010/website:v1 ."
docker build -t 192.168.1.12:5010/website:v1 .

echo "docker push 192.168.1.12:5010/website:v1"
docker push 192.168.1.12:5010/website:v1

# cd k8s
cd ..

echo "kubectl apply -f website.yml"
kubectl apply -f website.yml
