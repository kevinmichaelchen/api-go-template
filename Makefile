.PHONY: all
all:
	$(MAKE) gen-proto
	$(MAKE) gen-models

.PHONY: gen-proto
gen-proto:
	docker run --rm --volume "$(shell pwd):/workspace" --workdir /workspace bufbuild/buf mod update idl
	docker run --rm --volume "$(shell pwd):/workspace" --workdir /workspace bufbuild/buf lint
	docker run --rm --volume "$(shell pwd):/workspace" --workdir /workspace bufbuild/buf format -w
	docker run --rm --volume "$(shell pwd):/workspace" --workdir /workspace bufbuild/buf generate

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