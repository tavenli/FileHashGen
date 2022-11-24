go clean
set GOOS=windows
set GOARCH=386
go build -trimpath -ldflags "-w -s" -o FileHashGen.exe
