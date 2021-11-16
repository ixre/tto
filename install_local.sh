#!/usr/bin/env sh

export GOPROXY=https://goproxy.cn,direct
goods="windows"
if [ "$(uname)" = "Darwin" ]; then
  goods='darwin'
elif [ "$(uname)" = "Linux" ]; then
  goods="linux"
fi

CGO_ENABLED=0 GOOS=${goods} ARCH=amd64 go build -v -o tto bin/*.go
sudo rm -rf /usr/local/bin/tto
sudo mv tto /usr/local/bin

