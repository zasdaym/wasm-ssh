start:
	@docker compose up -d
	@GOOS=js GOARCH=wasm go build -o ./cmd/server/static/main.wasm ./cmd/client
	@go run ./cmd/server
