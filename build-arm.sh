#!/bin/bash
go clean
set GOOS=linux
set GOARCH=arm
go build -trimpath -ldflags "-w -s"
