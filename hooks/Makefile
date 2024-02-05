
.PHONY: prebuild
prebuild:
	rm -rf build
	mkdir -p build

.PHONY: deps
deps:
	npm install -g node-gyp
	npm i

.PHONY: lint
lint:
	npx tslint --project ./tsconfig.json --fix

.PHONY: test
test: build
	npm test

.PHONY: build
build: prebuild
	npx tsc

.PHONY: run
run:
	node --no-deprecation ./build/server.js hooks

.PHONY: reap
reap:
	node --no-deprecation ./build/server.js reap

.PHONY: publish
publish: publish-hooks publish-reap

.PHONY: publish-hooks
publish-hooks: test
	docker build -f Dockerfile.hooks -t registry.heroku.com/ttlsh-hooks/web .
	docker push registry.heroku.com/ttlsh-hooks/web
	heroku container:release web -a ttlsh-hooks

.PHONY: publish-reap
publish-reap: test
	docker build -f Dockerfile.reap -t registry.heroku.com/ttlsh-hooks/reap .
	docker push registry.heroku.com/ttlsh-hooks/reap
	heroku container:release reap -a ttlsh-hooks

