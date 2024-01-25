#!/bin/bash

bin_name=bootstrap

path=$1
function_name=$2

cd $path

rm -rf $bin_name $function_name.zip
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o $bin_name main.go
zip $function_name.zip $bin_name

cd ..