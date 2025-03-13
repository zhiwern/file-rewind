build:
	go-winres simply --icon .\assets\icon.ico && go build -o file-rewind.exe
	
run:
	go  run .\main.go

# Run once only
setup:
	go mod init file-rewind && go install github.com/tc-hib/go-winres@latest