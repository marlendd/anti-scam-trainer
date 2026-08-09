package config

import (
	"log"
	"os"
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
	SeedsPath      string        `yaml:"seeds_path" env:"SEEDS_PATH" env-default:"seeds"`
	// Connection pool
	MaxOpenConns    int           `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS" env-default:"50"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" env-default:"25"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME" env-default:"60m"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" env:"DB_CONN_MAX_IDLE_TIME" env-default:"5m"`

	SecureCookies bool `yaml:"secure_cookies" env:"SECURE_COOKIES" env-default:"false"`
	// CORS
	AllowedOrigins string `yaml:"allowed_origins" env:"ALLOWED_ORIGINS" env-default:"http://localhost:3000"`
	// SMTP
	SMTPHost     string `yaml:"smtp_host" env:"SMTP_HOST"`
	SMTPPort     int    `yaml:"smtp_port" env:"SMTP_PORT" env-default:"587"`
	SMTPUsername string `yaml:"smtp_username" env:"SMTP_USERNAME"`
	SMTPPassword string `yaml:"smtp_password" env:"SMTP_PASSWORD"`
	SMTPFrom     string `yaml:"smtp_from" env:"SMTP_FROM"`

	AppBaseURL string `yaml:"app_base_url" env:"APP_BASE_URL" env-default:"http://localhost:3000"`
}

func MustLoad(configPath string) Config {
	var cfg Config
	if _, err := os.Stat(configPath); err == nil {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Fatalf("cannot read config %q: %s", configPath, err)
		}
		return cfg
	}

	// файла нет — читаем конфигурацию только из переменных окружения
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config from environment: %s", err)
	}
	return cfg
}
