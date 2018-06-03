#!/bin/bash

mkdir -p build

### Linux ###
# (export GOOS=linux; export GOARCH=386; go build -o build/ng-deploy-$GOOS-$GOARCH ng-deploy.go)
(export GOOS=linux; export GOARCH=amd64; go build -o build/ng-deploy-$GOOS-$GOARCH ng-deploy.go)

### Windows ###
# (export GOOS=windows; export GOARCH=386; go build -o build/ng-deploy-$GOOS-$GOARCH.exe ng-deploy.go)
# (export GOOS=windows; export GOARCH=amd64; go build -o build/ng-deploy-$GOOS-$GOARCH.exe ng-deploy.go)