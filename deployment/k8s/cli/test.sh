#!/bin/zsh

echo "########## exec ############"
echo "./cli -name toJson_p -i input.txt -o ooooooo -l -j"
./cli -name toJson_p -i input.txt -o ooooooo -l -j


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
