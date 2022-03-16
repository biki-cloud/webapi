#!/bin/zsh


buildDir=./files
mkdir -p ${buildDir}

# ゲートウェイサーバは2台
for i in 1 2
do
  echo "apigw"${i}
  lbDir=${buildDir}/apigw${i}
  mkdir -p ${lbDir}
  rm -fr ${lbDir}/apigw
  go build -o ${lbDir}/apigw ../../microservices/apigw/cmd/application/main/main.go
done

# プログラムサーバは2台
for i in 1 2
do
  echo "server"${i}
  serverDir=${buildDir}/exec${i}
  mkdir -p ${serverDir}
  rm -fr ${serverDir}/exec
  rm -fr ${serverDir}/tmplates
  go build -o ${serverDir}/exec ../../microservices/exec/cmd/application/main/main.go
done


# websiteサーバは2台
for i in 1
do
  echo "website"${i}
  serverDir=${buildDir}/website${i}
  mkdir -p ${serverDir}
  rm -fr ${serverDir}/website
  rm -fr ${serverDir}/ui
  go build -o ${serverDir}/website ../../microservices/website/cmd/application/main/main.go
  cp -r ../../microservices/website/cmd/application/ui ${serverDir}/
done

echo "cli"
localDir=${buildDir}/cli
mkdir -p ${localDir}
go build -o ${localDir}/cli ../../cli/cmd/application/main/main.go

echo "process kill"
./server_kill.sh

echo "run all server"
./server_run.sh
