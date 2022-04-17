#!/bin/zsh

# current dir is k8s

echo -e "\n\n\n\n\n"
echo "################ config map ######################"
echo "kubectl delete -f configmap.yml"
kubectl delete -f configmap.yml

echo "kubectl apply -f configmap.yml"
kubectl apply -f configmap.yml

# 順番はこんな感じ
./exec-python1.sh
./exec-python2.sh

./apigw.sh

./website.sh


./cli.sh
