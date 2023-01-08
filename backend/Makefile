
build: FORCE
	@rm -rf _build
	go build -o _build/secretsanta ./cmd/secretsanta/main.go

test: FORCE
	go test -count=1 --race ./...

FORCE: