# envconf

Loads configuration structs from environment variables, with best-effort
`.env` loading for local development. A thin, opinionated wrapper around
[`caarlos0/env`](https://github.com/caarlos0/env) and
[`joho/godotenv`](https://github.com/joho/godotenv).

`envconf` is developed as part of the [gp-system](https://github.com/gp-system)
tooling but has no dependency on it: it works in any Go program that wants to
compose env-tagged config structs from independent packages into one typed
config.

## The problem it solves

Modular backends split configuration across packages: a server package owns
`LISTEN_ADDR`, a database package owns `DB_HOST`, an auth package owns
`JWT_SECRET`. Each package should own its own env keys without the
application gluing together dozens of `os.Getenv` calls or a hand-rolled
config loader. `envconf` composes those structs by embedding, resolves
`envPrefix` nesting, loads a `.env` file for local development without
requiring anything to be exported, and turns a missing required variable into
a single actionable error instead of a silent zero value.

## Install

```sh
go get github.com/gp-system/envconf@latest
```

Tagged releases start at `v0.1.0`.

## Usage

Compose config structs owned by different packages into one application
config:

```go
type Config struct {
	App AppConfig
	DB  DBConfig `envPrefix:"DB_"`
}

type AppConfig struct {
	ListenAddr string `env:"LISTEN_ADDR" envDefault:":8080"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"info"`
}

type DBConfig struct {
	Host string `env:"HOST,required"`
}

cfg := envconf.MustLoad[Config]()
```

### Load / MustLoad

`Load[T]` parses environment variables into `T`, loading a `.env` file in the
working directory first if one is present (existing environment variables are
never overridden, so a real deployment's env always wins over `.env`):

```go
cfg, err := envconf.Load[Config]()
if err != nil {
	log.Fatal(err)
}
```

`MustLoad[T]` is `Load` but panics on error. It is meant for `main()`, where
there is nothing more useful to do with a broken config than fail fast:

```go
cfg := envconf.MustLoad[Config]()
```

### LoadPrefixed

`LoadPrefixed[T]` parses only variables under a common prefix into `T`,
useful for loading the same struct type multiple times under different
prefixes (e.g. a second database connection):

```go
analytics, err := envconf.LoadPrefixed[DBConfig]("ANALYTICS_DB_")
```

## Design rules

- **Struct tags are the contract.** Every field uses the `env`/`envDefault`/
  `envPrefix` tags from `caarlos0/env`; `envconf` does not add its own tag
  vocabulary.
- **`.env` is best-effort.** A missing `.env` file is not an error: it only
  matters in local development, and a real deployment sets variables through
  its own mechanism.
- **Fail fast on missing required variables.** A required field left unset is
  a startup error (via `Load`) or a panic (via `MustLoad`), never a silent
  zero value.
