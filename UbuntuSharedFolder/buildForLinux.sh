@echo off
set GOOS=linux
set GOARCH=amd64

go build -o UbuntuSharedFolder/pglogalyze ./cmd/pglogalyze