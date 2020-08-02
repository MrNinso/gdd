#!/bin/bash

rm -rf ./build
mkdir ./build

GOOS=linux GOARCH=amd64 go build -o build/gdd.linux64 . &
pids[0]=$!
GOOS=linux GOARCH=amd64 go build -o build/gdd.linux32 . &
pids[1]=$!
GOOS=darwin GOARCH=amd64 go build -o build/gdd.mac64 . &
pids[2]=$!
GOOS=darwin GOARCH=amd64 go build -o build/gdd.mac32 . &
pids[3]=$!


for pid in ${pids[*]}; do
    wait $pid
done

