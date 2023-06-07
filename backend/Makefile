
build: FORCE
	@rm -rf _build
	go build -o _build/secretsanta ./cmd/secretsanta/main.go

test: FORCE
	go test -count=1 ./...

test.race: FORCE
	go test -count=1 --race ./...

dev.run: FORCE
	@go run ./cmd/devtools/main.go watcher

# svelte-check is the main way to enforce types in the frontend project. However, some of the errors are incorrect. So
# the project uses a custom checker that removes the useless error from the list.
dev.svelte-check: FORCE
	@go run ./cmd/devtools/main.go svelte-check

FORCE: