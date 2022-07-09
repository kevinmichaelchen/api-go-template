.PHONY: all
all:
	$(MAKE) gen-proto
	$(MAKE) gen-models

.PHONY: gen-proto
gen-proto:
	buf mod update idl
	buf lint idl
	buf format idl -w
	buf generate idl

.PHONY: gen-models
gen-models:
	sqlboiler psql --output internal/models

.PHONY: migrate-up
migrate-up:
	docker run -v $(shell pwd)/schema:/migrations \
	  --network host \
	  migrate/migrate \
	  -path=/migrations/ \
	  -database postgres://postgres:postgres@localhost:5432/foo\?sslmode=disable \
	  up

.PHONY: migrate-down
migrate-down:
	docker run -v $(shell pwd)/schema:/migrations \
	  --network host \
	  migrate/migrate \
	  -path=/migrations/ \
	  -database postgres://postgres:postgres@localhost:5432/foo\?sslmode=disable \
	  down -all