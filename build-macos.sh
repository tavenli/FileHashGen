#!/bin/bash
go clean
set GOOS=darwin
set GOARCH=amd64
go build -trimpath -ldflags "-w -s"
