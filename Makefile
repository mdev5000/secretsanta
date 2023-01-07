SPA_EMBED_PATH = "./backend/cmd/secretsanta/embedded"

schemas.refresh: schemas.clean schemas.gen

schemas.gen:
	npm run schema.gen

schemas.clean:
	npm run schema.clean

# Fully rebuild the project.
build.refresh: schemas.refresh build

build: build.frontend
	@rm -rf $(SPA_EMBED_PATH)
	@mkdir $(SPA_EMBED_PATH)
	# move the frontend assets into the backend
	@mv ./frontend/assets/ "$(SPA_EMBED_PATH)/assets"
	@mv ./frontend/build/ "$(SPA_EMBED_PATH)/spa"
	$(MAKE) -C ./backend build

build.frontend:
	$(MAKE) -C ./frontend build

FORCE: