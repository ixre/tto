#!/usr/bin/env sh

export GOPROXY=https://goproxy.io
CGO_ENABLED=0 GOOS=linux ARCH=amd64 go build -v -o tto bin/tto.go
sudo mv tto /usr/local/bin