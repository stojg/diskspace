build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dir_stats_linux_amd64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o dir_stats_darwin_amd64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o dir_stats_windows.exe .
