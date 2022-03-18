#!/bin/zsh

# current dir is k8s


echo "################ nfs server ######################"
echo "kubectl delete -f nfs-server.yml"
kubectl delete -f nfs-server.yml
echo "kubectl apply -f nfs-server.yml"
kubectl apply -f nfs-server.yml

echo "################ create volume with nfs ################"
#echo "kubectl delete -f nfs.yml"
#kubectl delete -f nfs.yml
echo "kubectl apply -f nfs.yml"
kubectl apply -f nfs.yml


echo "################ config map ######################"
echo "kubectl delete -f configmap.yml"
kubectl delete -f configmap.yml

echo "kubectl apply -f configmap.yml"
kubectl apply -f configmap.yml

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

cd deployment/k8s

echo "docker rmi website:v1"
docker rmi website:v1

cd website-kube

echo "docker build -t website:v1 ."
docker build -t website:v1 .

echo "minikube image rm website:v1"
minikube image rm website:v1

echo "minikube image load website:v1"
minikube image load website:v1

# cd k8s
cd ..

echo "kubectl apply -f website.yml"
kubectl apply -f website.yml



# ----------- apigw
echo "################### apigw #####################"

echo "kubectl delete -f apigw.yml"
kubectl delete -f apigw.yml


# cd to project root
cd ../..


echo "rm -fr deployment/k8s/apigw-kube/apigw/apigw"
rm -fr deployment/k8s/apigw-kube/apigw/apigw

echo "GOOS=linux GOARCH=amd64 go build -o deployment/k8s/apigw-kube/apigw/apigw microservices/apigw/cmd/application/main/main.go"
GOOS=linux GOARCH=amd64 go build -o deployment/k8s/apigw-kube/apigw/apigw microservices/apigw/cmd/application/main/main.go

cd deployment/k8s

echo "docker rmi apigw:v1"
docker rmi apigw:v1

cd apigw-kube

# build from Dockerfile
echo "docker build -t apigw:v1 ."
docker build -t apigw:v1 .

echo "minikube image rm apigw:v1"
minikube image rm apigw:v1


echo "minikube image load apigw:v1"
minikube image load apigw:v1

# cd k8s
cd ..



echo "kubectl apply -f apigw.yml"
kubectl apply -f apigw.yml




# ----------- exec1
echo "######################### exec1 ###########################"

echo "kubectl delete -f exec1.yml"
kubectl delete -f exec1.yml


# cd to project root
cd ../..


echo "rm -fr deployment/k8s/exec-kube1/exec/exec"
rm -fr deployment/k8s/exec-kube1/exec/exec

echo "GOOS=linux GOARCH=amd64 go build -o deployment/k8s/exec-kube1/exec/exec microservices/exec/cmd/application/main/main.go"
GOOS=linux GOARCH=amd64 go build -o deployment/k8s/exec-kube1/exec/exec microservices/exec/cmd/application/main/main.go

cd deployment/k8s

echo "docker rmi exec1:v1"
docker rmi exec1:v1

cd exec-kube1

# build from Dockerfile
echo "docker build -t exec1:v1 ."
docker build -t exec1:v1 .

echo "minikube image rm exec1:v1"
minikube image rm exec1:v1


echo "minikube image load exec1:v1"
minikube image load exec1:v1

# cd k8s
cd ..

echo "kubectl apply -f exec1.yml"
kubectl apply -f exec1.yml



# ----------- exec2

echo "################### exec2 ########################"

echo "kubectl delete -f exec2.yml"
kubectl delete -f exec2.yml


# cd to project root
cd ../..


echo "rm -fr deployment/k8s/exec-kube2/exec/exec"
rm -fr deployment/k8s/exec-kube2/exec/exec

echo "GOOS=linux GOARCH=amd64 go build -o deployment/k8s/exec-kube2/exec/exec microservices/exec/cmd/application/main/main.go"
GOOS=linux GOARCH=amd64 go build -o deployment/k8s/exec-kube2/exec/exec microservices/exec/cmd/application/main/main.go

cd deployment/k8s

echo "docker rmi exec2:v1"
docker rmi exec2:v1

cd exec-kube2

# build from Dockerfile
echo "docker build -t exec2:v1 ."
docker build -t exec2:v1 .

echo "minikube image rm exec2:v1"
minikube image rm exec2:v1

echo "minikube image load exec2:v1"
minikube image load exec2:v1

# cd k8s
cd ..

echo "kubectl apply -f exec2.yml"
kubectl apply -f exec2.yml


echo "##################### cli ############################"
mkdir -p cli
rm -fr cli/cli
go build -o cli/cli ../../cli/cmd/application/main/main.go
