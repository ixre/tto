#!/usr/bin/env sh

CGO_ENABLED=0 GOOS=linux ARCH=amd64 go build -o tto bin/tto.go
sudo mv tto /usr/local/bin