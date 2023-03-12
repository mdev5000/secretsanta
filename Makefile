SPA_EMBED_PATH = "./backend/cmd/secretsanta/embedded"

deps.install: FORCE
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/google/go-jsonnet/cmd/jsonnet@latest

schemas.gen: schemas.gen.backend schemas.gen.frontend

schemas.gen.backend:
	protoc -I=protos --go_out=backend/internal protos/*

schemas.gen.frontend:
	npx protoc --ts_out frontend/src/lib/requests --proto_path protos protos/*

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

dev.run: FORCE
	$(MAKE) -C ./backend dev.run

uitest.build:
	docker-compose -f ./docker/dc-uitesting.yml build

uitest.up:
	docker-compose -f ./docker/dc-uitesting.yml up -d

uitest.down:
	docker-compose -f ./docker/dc-uitesting.yml down

FORCE: