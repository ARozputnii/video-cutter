tiny:
	go mod tidy

run:
	go run cmd/main.go

windows-build: clean-windows-exe
	GOOS=windows GOARCH=amd64 go build -o video-cutter.exe cmd/main.go

clean-windows-exe:
	rm -f video-cutter.exe

clean-uploads:
	rm -rf uploads
