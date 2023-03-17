
sp:
	npm run smui-prepare

install :FORCE
	npm install

build: FORCE
	npm run build
	$(MAKE) assets.replacePaths

check: FORCE
	npm run check

assets.replacePaths: FORCE
	@echo "Fixing asset paths..."
	@find ./assets -type f -name "*.js" -exec sed -i 's|http://localhost:3000/replaceme/assets|/assets|g' {} \;
	@find ./build -type f -exec sed -i 's|http://localhost:3000/replaceme/assets|/assets|g' {} \;
	@echo "Asset paths fixed."

test: FORCE
	# run unit tests
	npx vitest --run

FORCE: