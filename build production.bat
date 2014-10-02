@echo off

go run build.go production

if errorlevel 1 pause