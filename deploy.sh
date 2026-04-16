#!/usr/bin/env bash

#rm -rf output
#go install cmd/*.go
go mod tidy

CGO_ENABLED=0 GOOS=linux ARCH=amd64 go build -o tto cmd/*.go
CGO_ENABLED=0 GOOS=darwin ARCH=amd64 go build -o mac-darwin-amd64 cmd/*.go
CGO_ENABLED=0 GOOS=darwin ARCH=amd64 go build -o mac-darwin-arm64 cmd/*.go
CGO_ENABLED=0 GOOS=windows ARCH=amd64 go build -o tto.exe cmd/*.go

tar cvzf tto-generator-bin.tar.gz \
  generate.sh mac-darwin-amd64 mac-darwin-arm64 tto \
  tto.exe tto.conf templates README.md LICENSE \
  
#rm -rf tto mac-tto tto.exe
