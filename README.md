# api-go-template

[![Lines Of Code](https://tokei.rs/b1/github/kevinmichaelchen/api-go-template?category=code)](https://github.com/kevinmichaelchen/api-go-template)

A boilerplate Go repo that comes with:

* **Dependency Injection** using [fx](https://github.com/uber-go/fx)
* **Env Var Configs** using [go-envconfig](https://github.com/sethvargo/go-envconfig)
* **gRPC** using [connect-go](https://github.com/bufbuild/connect-go)
* **DB Migration Tool** using [migrate](https://github.com/golang-migrate/migrate)
* **ORM** using [sqlboiler](https://github.com/volatiletech/sqlboiler)
* **Command-Line Interface** using [cobra](https://github.com/spf13/cobra)

## Project structure

| Directory                                    | Description                               |
|----------------------------------------------|-------------------------------------------|
| [`./cmd`](./cmd)                             | CLI for making gRPC requests              |
| [`./idl`](./idl)                             | Protobufs (Interface Definition Language) |
| [`./internal/app`](./internal/app)           | App dependency injection / initialization |
| [`./internal/idl`](./internal/idl)           | Auto-generated protobufs                  |
| [`./internal/models`](./internal/models)     | Auto-generated ORM / models               |
| [`./internal/service`](./internal/service)   | Service layer / Business logic            |
| [`./schema`](./schema)                       | SQL migration scripts                     |

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

## Database
### Run Migrations
```bash
make migrate-up
```
or you can run:
```bash
docker run -v $(pwd)/schema:/migrations \
  --network host \
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