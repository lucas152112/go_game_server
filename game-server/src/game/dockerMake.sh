#!/bin/bash
# GONOPROXY=https://goproxy.io,direct
#docker run --rm  \
#   -v "$(dirname $(dirname $PWD))":/poker \
#   -w /poker/src/game/  -e GOPATH=/poker/ -e GONOPROXY="https://goproxy.cn,direct" \
#   golang:1.15.2  go mod tidy && \
#    go mod download && \
#    go mod vendor && \
#    go build -v -o /poker/src/script/game
#


docker run --rm  \
   -v "$(dirname $(dirname $PWD))":/poker \
   -w /poker/src/game/  -e GOPATH=/poker/ -e GONOPROXY="https://goproxy.cn,direct" \
   golang:1.15.2  go build -v -o /poker/src/script/game