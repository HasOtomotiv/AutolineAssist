env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-w -s" -o ./bin/AutolineAssist.exe ./cmd
