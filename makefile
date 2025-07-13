port=8080
.phony: run
run: 
	go run ./cmd/app/main.go
stop:
	@fuser -k $(port)/tcp || true
restart: stop run