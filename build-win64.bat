go clean
set GOOS=windows
set GOARCH=amd64
go build -trimpath -ldflags "-w -s" -o FileHashGen-x64.exe
