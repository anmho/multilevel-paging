

all: project2



project2:
	go build -o ./bin/project2 main.go

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o ./bin/project2-linux-amd64 main.go

darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o ./bin/project2-darwin-arm64 main.go
