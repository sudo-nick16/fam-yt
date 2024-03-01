web:
	cd web && pnpm run dev

server: s
	go run cmd/server/main.go

fetcher: f
	go run cmd/fetcher/main.go

build-server: s
	go build -o bin/server cmd/server/main.go

build-fetcher: f
	go build -o bin/fetcher cmd/fetcher/main.go

s: cmd/server/main.go
f: cmd/fetcher/main.go
