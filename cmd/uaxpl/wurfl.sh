#!/bin/bash

echo "compile binary"
go build .
echo "binary compiled successfully"

echo "build diff for wurfl dataset #0"
./uaxpl -url=https://github.com/koykov/dataset/raw/refs/heads/master/ua/wurfl0.json --threads=8 --out=wurfl0.diff.txt --verbose
echo "wurfl dataset #0 done"

echo "build diff for wurfl dataset #1"
./uaxpl -url=https://github.com/koykov/dataset/raw/refs/heads/master/ua/wurfl1.json --threads=8 --out=wurfl1.diff.txt --verbose
echo "wurfl dataset #1 done"
