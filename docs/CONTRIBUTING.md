## Lines of code
```bash
tokei
```

## Protobufs
To regenerate protobufs, run 
```bash
make gen-proto
```

### Running migrations

To run database migrations:

```bash
docker run -v $(pwd)/schema:/migrations \
  --network host migrate/migrate \
  -path=/migrations/ \
  -database postgres://postgres:postgres@localhost:5432/foo\?sslmode=disable \
  up
```

### Generating SQLBoiler code

We use [sqlboiler](https://github.com/volatiletech/sqlboiler) to auto-generate
a strongly-typed ORM by pointing it at our current schema.

```bash
# Generate code
make gen-models
```

## Connecting to postgres
```bash
psql postgres://postgres:postgres@localhost:5432/foo
```