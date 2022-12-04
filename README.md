# api-go-template

[![Lines Of Code](https://aschey.tech/tokei/github/kevinmichaelchen/api-go-template?category=code&style=for-the-badge)](https://github.com/kevinmichaelchen/api-go-template)
![](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)

A boilerplate Go repo that comes with:

* **Dependency Injection** using [uber-go/fx](https://github.com/uber-go/fx) (:tv: [video](https://www.youtube.com/watch?v=nLskCRJOdxM))
* **DB Migration Tool** using [golang-migrate](https://github.com/golang-migrate/migrate) (:tv: [video](https://youtu.be/ZRUEJX1fqYc?t=845))
* **ORM** using [sqlboiler](https://github.com/volatiletech/sqlboiler) (:tv: [video](https://www.youtube.com/watch?v=M9bgMOLQLs8))
* **Env Var Configs** using [go-envconfig](https://github.com/sethvargo/go-envconfig)
* **gRPC** using [connect-go](https://github.com/bufbuild/connect-go)
* **Protobufs**, compiled, formatted, linted, and more with [Buf](https://buf.build/)
* **Command-Line Interface** using [cobra](https://github.com/spf13/cobra)

These libraries do a lot of heavy lifting in terms of boilerplate.

For example:
* The **fx** framework manages dependency injection and application life-cycle 
for you. (See [`example_test.go`](https://github.com/uber-go/fx/blob/master/example_test.go)).
* **sqlboiler** makes DB CRUD really simple. (See [`foo.go`](https://github.com/kevinmichaelchen/api-go-template/blob/main/internal/service/db/foo.go)).
[Raw queries](https://github.com/volatiletech/sqlboiler#raw-query) are supported
as an escape hatch.

## Project structure

| Directory                                        | Description                               |
|--------------------------------------------------|-------------------------------------------|
| [`./cmd`](./cmd)                                 | CLI for making gRPC requests              |
| [`./idl`](./idl)                                 | Protobufs (Interface Definition Language) |
| [`./internal/app`](./internal/app)               | App dependency injection / initialization |
| [`./internal/idl`](./internal/idl)               | Auto-generated protobufs                  |
| [`./internal/models`](./internal/models)         | Auto-generated ORM / models               |
| [`./internal/service`](./internal/service)       | Service layer / Business logic            |
| [`./internal/service/db`](./internal/service/db) | Data layer                                |
| [`./schema`](./schema)                           | SQL migration scripts                     |

## Getting started
```bash
# Step 1: Start containers in detached (background) mode
docker-compose up -d

# Step 2: Create the database schema
make migrate-up

# Step 3: Start app
go run main.go
```

Finally, hit the API (using [HTTPie](https://httpie.io/))
```bash
# Create a new Foo
http POST \
  http://localhost:8081/coop.drivers.foo.v1beta1.FooService/CreateFoo \
    name="Kevin"

# Get existing Foo
http POST \
  http://localhost:8081/coop.drivers.foo.v1beta1.FooService/GetFoo \
    id="cb4c4rnrirfucgsert7g"
```

Or with curl:
```bash
curl -X POST http://localhost:8081/coop.drivers.foo.v1beta1.FooService/CreateFoo \
  -H "Content-Type: application/json" \
  -d '{"name": "Kevin"}'

curl -X POST http://localhost:8081/coop.drivers.foo.v1beta1.FooService/GetFoo \
  -H "Content-Type: application/json" \
  -d '{"id": "cb4c4rnrirfucgsert7g"}'
```

## Database
### Run Migrations
```bash
make migrate-up
```
or you can run:
```bash
docker run -v $(pwd)/schema:/migrations \
  --network host \
  --rm \
  migrate/migrate \
  -path=/migrations/ \
  -database postgres://postgres:postgres@localhost:5432/foo\?sslmode=disable \
  up
```

This will run all migrations in [`./schema`](./schema).

### Create New Migration
To create a new migration called `create-new-table`, run:
```bash
docker run -v $(pwd)/schema:/migrations \
  --network host \
  --rm \
  migrate/migrate \
  -path=/migrations/ \
  create \
  -dir /migrations \
  -ext sql \
  create-new-table
```

This will create a new _up_ and _down_ migration in [`./schema`](./schema).

### Auto-generate ORM DB Models
We have a sqlboiler command that introspects the DB and generates ORM models.
```bash
make gen-models
```
Configuration comes from [`sqlboiler.toml`](./sqlboiler.toml)

## Telemetry
### Traces
See traces in [Jaeger](https://www.jaegertracing.io/) [here](http://localhost:16686)

### Metrics
See metrics in [Prometheus](https://prometheus.io/) [here](http://localhost:9090/graph?g0.expr=_coop_drivers_foo_v1beta1_FooService_CreateFoo&g0.tab=1&g0.stacked=0&g0.show_exemplars=0&g0.range_input=15m).