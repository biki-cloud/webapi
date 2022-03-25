#!/bin/zsh

echo "########## exec ############"

echo "export APIGW_SERVER_URIS=http://192.168.59.102:30006"
export APIGW_SERVER_URIS=http://192.168.59.102:30006

echo "./cli -name toJson -i input.txt -o ooooooo -l -j"
./cli -name toJson -i input.txt -o ooooooo -l -j


echo ""
echo ""
echo ""
echo "########## result ############"
echo "cat ooooooo/input.txt.json"
cat ooooooo/input.txt.json

echo ""
echo ""
echo ""
echo "rm -fr ooooooo"
rm -fr ooooooo
