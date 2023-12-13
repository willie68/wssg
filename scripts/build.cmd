@echo off
echo building generator
go build -ldflags="-s -w" -o wssg.exe cmd/wssg/main.go