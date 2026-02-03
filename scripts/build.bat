@echo off
setlocal

set BINARY=portflow
set BUILD_DIR=bin

if not exist %BUILD_DIR% mkdir %BUILD_DIR%

echo Building for windows/amd64...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags "-s -w" -o %BUILD_DIR%/%BINARY%-windows-amd64.exe ./cmd/portflow

echo Building for linux/amd64...
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags "-s -w" -o %BUILD_DIR%/%BINARY%-linux-amd64 ./cmd/portflow

echo Build complete.
pause
