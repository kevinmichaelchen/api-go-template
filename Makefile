DOCKER_BUF_FLAGS = --rm --volume "$(shell pwd):/workspace" --workdir /workspace
# Buf CLI versions:
# https://hub.docker.com/r/bufbuild/buf/tags
DOCKER_BUF = bufbuild/buf:1.9.0

.PHONY: all
all:
	$(MAKE) buf-gen
	$(MAKE) gen-models

.PHONY: buf-lint
buf-lint:
	docker run $(DOCKER_BUF_FLAGS) $(DOCKER_BUF) lint
	docker run $(DOCKER_BUF_FLAGS) $(DOCKER_BUF) format -w
	docker run $(DOCKER_BUF_FLAGS) $(DOCKER_BUF) breaking --against 'https://github.com/kevinmichaelchen/api-go-template.git#branch=main'

.PHONY: buf-mod-update
buf-mod-update:
	@for i in $(shell fd buf.yaml | xargs dirname) ; do \
	  docker run $(DOCKER_BUF_FLAGS) $(DOCKER_BUF) mod update $$i ; \
	done

.PHONY: buf-gen
buf-gen:
	docker run $(DOCKER_BUF_FLAGS) $(DOCKER_BUF) generate

.PHONY: gen-models
gen-models:
	sqlboiler psql --output internal/models

.PHONY: migrate-up
migrate-up:
	docker run -v $(shell pwd)/schema:/migrations \
	  --rm \
	  --network host \
	  migrate/migrate \
	  -path=/migrations/ \
	  -database postgres://postgres:postgres@localhost:5432/foo\?sslmode=disable \
	  up

.PHONY: migrate-down
migrate-down:
	docker run -v $(shell pwd)/schema:/migrations \
	  --rm \
	  --network host \
	  migrate/migrate \
	  -path=/migrations/ \
	  -database postgres://postgres:postgres@localhost:5432/foo\?sslmode=disable \
	  down -all