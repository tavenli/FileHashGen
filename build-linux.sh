#!/bin/bash
go clean
set GOOS=linux
set GOARCH=amd64
go build -trimpath -ldflags "-w -s"
