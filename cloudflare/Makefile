
.PHONY: deps
deps:
	yarn global add node-gyp
	yarn --silent --frozen-lockfile

.PHONY: prebuild
prebuild:
	rm -rf build
	mkdir -p build

.PHONY: lint
lint:
	npx tslint --project ./tsconfig.json --fix

.PHONY: build
build: prebuild
	`yarn bin`/webpack

.PHONY: preview
preview: build
	`yarn bin`/workers-preview < dist/bundle.js
