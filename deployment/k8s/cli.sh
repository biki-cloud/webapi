#!/bin/zsh



echo -e "\n\n\n\n\n"
echo "##################### cli ############################"
echo "mkdir -p cli"
mkdir -p cli

echo "rm -fr cli/cli"
rm -fr cli/cli

echo "go build -o cli/cli ../../cli/cmd/application/main/main.go"
go build -o cli/cli ../../cli/cmd/application/main/main.go
