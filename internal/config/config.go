package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	EnvLocal       = "local"
	EnvDevelopment = "dev"
	EnvProduction  = "prod"
)

type Bot struct {
	Token string `yaml:"token" env-required:"true"`
}

type Marzban struct {
	BaseURL string `yaml:"base_url"`
	Token   string `yaml:"token"`
	Enabled bool   `yaml:"enabled" env-default:"false"`
}

type Postgres struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

type Config struct {
	Env      string   `yaml:"env" env-required:"true"`
	Bot      Bot      `yaml:"bot" env-required:"true"`
	Postgres Postgres `yaml:"postgres" env-required:"true"`
	AdminIDs []int64  `yaml:"admin_ids"`
}

func (p Postgres) DSN() string {
	return "postgres://" + p.User + ":" + p.Password + "@" + p.Host + ":" + p.Port + "/" + p.Name + "?sslmode=" + p.SSLMode
}

func MustLoad() *Config {
	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("error to get config path")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config 
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func (c *Config) IsProd() bool {
	return c.Env==EnvProduction
}

func (c *Config) IsAdmin(telegramID int64) bool {
	for _, id := range c.AdminIDs {
		if id == telegramID {
			return true
		}
	}
	return false
}