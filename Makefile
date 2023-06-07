SPA_EMBED_PATH = "./backend/cmd/secretsanta/embedded"
GO_PROTO_PATH_IN = "./backend/_build/goprotos"
GO_PROTO_PATH_OUT = "./backend/internal/requests/gen"

# dependencies ---------------------------------------------------------------------------------------------------------

deps.install: FORCE
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/google/go-jsonnet/cmd/jsonnet@latest

# updates --------------------------------------------------------------------------------------------------------------

updates.list: FORCE
	./scripts/updates-list.sh

updates.install: FORCE
	./scripts/updates-install.sh

# protobuf schema ------------------------------------------------------------------------------------------------------

schemas.gen: schemas.gen.backend schemas.gen.frontend

schemas.gen.backend:
	@rm -rf $(GO_PROTO_PATH_IN)
	@mkdir -p $(GO_PROTO_PATH_IN)
	protoc -I=protos --go_out=$(GO_PROTO_PATH_IN) protos/core/*
	protoc -I=protos --go_out=$(GO_PROTO_PATH_IN) protos/setup/*
	@rm -rf $(GO_PROTO_PATH_OUT)
	@mv $(GO_PROTO_PATH_IN)/github.com/mdev5000/secretsanta/internal/requests/gen $(GO_PROTO_PATH_OUT)

schemas.gen.frontend:
	npx protoc --ts_out frontend/src/lib/requests --proto_path protos protos/core/*
	npx protoc --ts_out frontend/src/lib/requests --proto_path protos protos/setup/*


# building -------------------------------------------------------------------------------------------------------------

build.ci: deps.install frontend.install build.refresh

# Fully rebuild the project.
build.refresh: schemas.gen build

# Compiles the frontend SPA and moves the resulting files into the backend project to be embedded in the resulting
# backend binary and then compiles the backend binary.
build: build.frontend
	@rm -rf $(SPA_EMBED_PATH)
	@mkdir $(SPA_EMBED_PATH)
	# move the frontend assets into the backend
	@mv ./frontend/assets/ "$(SPA_EMBED_PATH)/assets"
	@mv ./frontend/build/ "$(SPA_EMBED_PATH)/spa"
	$(MAKE) -C ./backend build

frontend.install:
	$(MAKE) -C ./frontend install

build.frontend:
	$(MAKE) -C ./frontend build

# dev environment ------------------------------------------------------------------------------------------------------

dev.docker.up: FORCE
	docker-compose up -d

dev.docker.down: FORCE
	docker-compose down

dev.run: FORCE
	$(MAKE) -C ./backend dev.run

# testing --------------------------------------------------------------------------------------------------------------

# test everything locally
test.all:
	# Make sure we are not running development docker
	$(MAKE) dev.docker.down
# @todo add front-end check back in once fixed.
	@echo "Check frontend typing..."
	$(MAKE) -C ./backend dev.svelte-check
	@echo "Check frontend typing (done)"
	@echo "Test frontend..."
	$(MAKE) -C ./frontend test
	@echo "Test frontend (done)."
	@echo "Test backend..."
	$(MAKE) -C ./backend test.race
	@echo "Test backend (done)."
	@echo "UI testing..."
	$(MAKE) uitest.test.build
	@echo "UI testing (done)."
	@echo ""
	@echo "All testing passed! :)"

# ui testing -----------------------------------------------------------------------------------------------------------

uitest.test.rebuild: uitest.down uitest.test.build

uitest.test.build: uitest.build
	$(MAKE) uitest.test

uitest.test:
	./scripts/run-test-ui.sh

uitest.build:
	docker-compose -f ./docker/dc-uitesting.yml build

uitest.up:
	docker-compose -f ./docker/dc-uitesting.yml up -d

uitest.down:
	docker-compose -f ./docker/dc-uitesting.yml down

FORCE: