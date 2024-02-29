server:
	go run cmd/server/main.go

fetcher:
	go run cmd/fetcher/main.go

build-server: 
	go build -o bin/server cmd/server/main.go

build-fetcher: 
	go build -o bin/fetcher cmd/fetcher/main.go

server: cmd/server/main.go
fetcher: cmd/fetcher/main.go
