// Package envconf loads configuration structs from environment variables
// using caarlos0/env, with best-effort .env loading for development.
//
// Independent packages can each own an env-tagged config struct;
// applications compose them with envPrefix:
//
//	type Config struct {
//		App AppConfig
//		DB  DBConfig `envPrefix:"DB_"`
//	}
//
//	type AppConfig struct {
//		ListenAddr string `env:"LISTEN_ADDR" envDefault:":8080"`
//		LogLevel   string `env:"LOG_LEVEL" envDefault:"info"`
//	}
//
//	type DBConfig struct {
//		Host string `env:"HOST,required"`
//	}
//
//	cfg := envconf.MustLoad[Config]()
package envconf

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Load parses environment variables into T. A .env file in the working
// directory is loaded first if present (existing variables are never
// overridden), so local development works without exporting anything.
func Load[T any]() (T, error) {
	_ = godotenv.Load()
	cfg, err := env.ParseAs[T]()
	if err != nil {
		return cfg, fmt.Errorf("envconf: %w", err)
	}
	return cfg, nil
}

// MustLoad is Load but panics on error. Intended for main().
func MustLoad[T any]() T {
	cfg, err := Load[T]()
	if err != nil {
		panic(err)
	}
	return cfg
}

// LoadPrefixed parses environment variables that share a common prefix into
// T, e.g. LoadPrefixed[pg.Config]("NEWS_DB_").
func LoadPrefixed[T any](prefix string) (T, error) {
	_ = godotenv.Load()
	var cfg T
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: prefix}); err != nil {
		return cfg, fmt.Errorf("envconf: prefix %q: %w", prefix, err)
	}
	return cfg, nil
}
