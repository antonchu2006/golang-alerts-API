@echo off
go build -ldflags "-s -w" alert_server.go 
echo Compiled