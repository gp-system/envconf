package envconf_test

import (
	"testing"
	"time"

	"github.com/gp-system/envconf"
)

// serverConfig stands in for a kit package's env-tagged config struct (e.g.
// gpsystem's server.Config) to exercise composition without depending on the
// kit itself: envconf is a leaf package with zero gp-system/gpsystem imports.
type serverConfig struct {
	ListenAddr      string        `env:"LISTEN_ADDR" envDefault:":8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
}

type dbConfig struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"PORT" envDefault:"5432"`
}

type appConfig struct {
	Server serverConfig
	DB     dbConfig `envPrefix:"DB_"`
	Secret string   `env:"APP_SECRET,required"`
}

func TestComposedLoad(t *testing.T) {
	t.Setenv("LISTEN_ADDR", ":4000")
	t.Setenv("SHUTDOWN_TIMEOUT", "3s")
	t.Setenv("DB_HOST", "db.internal")
	t.Setenv("APP_SECRET", "s3cr3t")

	cfg, err := envconf.Load[appConfig]()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.ListenAddr != ":4000" || cfg.Server.ShutdownTimeout != 3*time.Second {
		t.Errorf("server config not loaded: %+v", cfg.Server)
	}
	if cfg.DB.Host != "db.internal" || cfg.DB.Port != 5432 {
		t.Errorf("prefixed db config wrong: %+v", cfg.DB)
	}
	if cfg.Secret != "s3cr3t" {
		t.Errorf("secret = %q", cfg.Secret)
	}
}

func TestRequiredMissingFails(t *testing.T) {
	if _, err := envconf.Load[appConfig](); err == nil {
		t.Fatal("want error for missing APP_SECRET")
	}
}

func TestLoadPrefixed(t *testing.T) {
	t.Setenv("NEWS_DB_HOST", "h1")
	cfg, err := envconf.LoadPrefixed[dbConfig]("NEWS_DB_")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Host != "h1" {
		t.Errorf("host = %q", cfg.Host)
	}
}
