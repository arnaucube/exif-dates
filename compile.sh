#!/bin/sh

mkdir -p bin

echo "building linux binaries"
GOOS=linux GOARCH=amd64 go build -o bin/exif-dates-amd64-linux *.go

echo "building windows binaries"
GOOS=windows GOARCH=amd64 go build -o bin/exif-dates-amd64.exe *.go

echo "building macOS binaries"
GOOS=darwin GOARCH=amd64 go build -o bin/exif-dates-amd64-darwin *.go
GOOS=darwin GOARCH=arm64 go build -o bin/exif-dates-arm64-darwin *.go
