package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string        `yaml:"env" env:"ENV" env-default:"local"`
	Port           string        `yaml:"port" env:"PORT" env-default:"8080"`
	Timeout        time.Duration `yaml:"timeout" env:"API_TIMEOUT" env-default:"5s"`
	LogLevel       string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	DatabaseURL    string        `yaml:"database_url" env:"DATABASE_URL" env-required:"true"`
	MigrationsPath string        `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-default:"migrations"`
	// Connection pool
	MaxOpenConns    int           `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS" env-default:"50"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" env-default:"25"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME" env-default:"60m"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" env:"DB_CONN_MAX_IDLE_TIME" env-default:"5m"`

	SecureCookies bool `yaml:"secure_cookies" env:"SECURE_COOKIES" env-default:"false"`
}

func MustLoad(configPath string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %q: %s", configPath, err)
	}
	return cfg
}
