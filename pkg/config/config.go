package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type REST struct {
	Port string `required:"true" default:"5000"`
	Host string `required:"true" default:"localhost"`
}

func (g REST) Addr() string {
	return g.Host + ":" + g.Port
}

type Server struct {
	REST        REST
	LogLevel    string        `split_words:"true" required:"true" default:"info"`
	GracePeriod time.Duration `split_words:"true" required:"true" default:"60s"`
	MaxCPU      int           `split_words:"true" required:"true" default:"1"`
	Debug       bool          `required:"true" default:"false"`
}

type Database struct {
	Name            string        `required:"true" default:"portfolio"`
	Host            string        `required:"true" default:"localhost"`
	Port            string        `required:"true" default:"5432"`
	User            string        `required:"true" default:"root"`
	Password        string        `required:"true" default:"password"`
	ConnMaxLifetime time.Duration `split_words:"true" required:"true" default:"10m"`
	MaxOpenConns    int           `split_words:"true" required:"true" default:"30"`
	MaxIdleConns    int           `split_words:"true" required:"true" default:"30"`
	MigrationPath   string        `split_words:"true" required:"true"`
}

func (d *Database) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Name,
	)
}

type Config struct {
	Server   Server
	Database Database
}

func (c *Config) DBString() string {
	return c.Database.DSN()
}

func Load() (cfg Config, err error) {
	err = envconfig.Process("INVESTMENT", &cfg)
	return cfg, errors.Wrap(err, "cannot parse env configuration")
}
