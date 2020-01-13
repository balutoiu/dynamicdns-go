#!/bin/sh

echo "building docker images for ${GOOS}/${GOARCH} ..."

REPO="github.com/alinbalutoiu/dynamicdns-go"

export CGO_ENABLED=0
go build -o release/linux/${GOARCH}/dynamicdns-go ${REPO}
