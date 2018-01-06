package pylint

import (
	"github.com/caarlos0/env"
	"log"
)

// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

type HookResponse struct {
	Msg string
}

type Config struct {
	Github struct {
		// can be 0 to ignore admin
		AdminId       int64  `env:"PYLINTGO_GITHUB_ADMINID"         envDefault:"0"`
		IntegrationID int64  `env:"PYLINTGO_GITHUB_INTEGRATIONID"   envDefault:"0"`
		Secret        string `env:"PYLINTGO_GITHUB_SECRET"          envDefault:"dummy"`
		KeyPath       string `env:"PYLINTGO_GITHUB_KEYPATH"         envDefault:"/keys/key.pem"`
	}
	Database struct {
		Host     string `env:"PYLINTGO_DB_HOST"         envDefault:"postgres"`
		User     string `env:"PYLINTGO_DB_USER"         envDefault:"postgres"`
		Name     string `env:"PYLINTGO_DB_NAME"         envDefault:"postgres"`
		Password string `env:"PYLINTGO_DB_PASSWORD"     envDefault:"postgres"`
		Port     string `env:"PYLINTGO_DB_PORT"         envDefault:"5432"`
	}
	Redis struct {
		Host string `env:"PYLINTGO_REDIS_HOST"      envDefault:"pool"`
		Port string `env:"PYLINTGO_REDIS_PORT"      envDefault:"6379"`
	}
	Port       int    `env:"PYLINTGO_PORT"        envDefault:"4444"`
	PublicPort int    `env:"PYLINTGO_PUBLICPORT"  envDefault:"9097"`
	Url        string `env:"PYLINTGO_STATEURL"    envDefault:"http://subdomain.domain.com/"`
	Name       string `env:"PYLINTGO_NAME"        envDefault:"PyLinter"`
}

var Cfg Config

func (cfg *Config) Parse() {
	if err := env.Parse(cfg); err != nil {
		log.Fatal("Unable to parse config: ", err)
	}
	if err := env.Parse(&cfg.Github); err != nil {
		log.Fatal("Unable to parse config.GitHub: ", err)
	}
	if err := env.Parse(&cfg.Database); err != nil {
		log.Fatal("Unable to parse config.Database: ", err)
	}
	if err := env.Parse(&cfg.Redis); err != nil {
		log.Fatal("Unable to parse config.Redis: ", err)
	}
}
