#!/usr/bin/env bash

rm -rf output
go install bin/tto.go

CGO_ENABLED=0 GOOS=linux ARCH=amd64 go build -o tto bin/tto.go
CGO_ENABLED=0 GOOS=darwin ARCH=amd64 go build -o mac-tto bin/tto.go
CGO_ENABLED=0 GOOS=windows ARCH=amd64 go build -o tto.exe bin/tto.go

tar cvzf tto-build-bin.tar.gz tto.sh mac-tto tto\
 tto.exe tto.conf templates README.md
rm -rf tto mac-tto tto.exe
