web:
	cd web && pnpm run dev

server: keys
	go run cmd/server/main.go

fetcher:
	go run cmd/fetcher/main.go

keys: cert.pem key.pem

cert.pem key.pem:
	./scripts/generate-keys.sh

build-server: 
	go build -o bin/server cmd/server/main.go

build-fetcher: 
	go build -o bin/fetcher cmd/fetcher/main.go

server: cmd/server/main.go
fetcher: cmd/fetcher/main.go
