@echo off
echo building generator
go build -ldflags="-s -w" -o wssg.exe main.go