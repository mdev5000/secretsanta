SPA_EMBED_PATH = "./backend/cmd/secretsanta/embedded"

deps.install: FORCE
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/google/go-jsonnet/cmd/jsonnet@latest

schemas.gen: schemas.gen.backend schemas.gen.frontend

schemas.gen.backend:
	protoc -I=protos --go_out=backend/internal protos/*

schemas.gen.frontend:
	npx protoc --ts_out frontend/src/lib/requests --proto_path protos protos/*

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

build.frontend:
	$(MAKE) -C ./frontend build

dev.run: FORCE
	$(MAKE) -C ./backend dev.run

FORCE: