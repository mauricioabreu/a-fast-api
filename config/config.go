package config

import (
	"log"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	DBUser          string `env:"DB_USER,required"`
	DBPassword      string `env:"DB_PASSWORD,required"`
	DBName          string `env:"DB_NAME,required"`
	DBHost          string `env:"DB_HOST,required"`
	DBPort          string `env:"DB_PORT" envDefault:"5432"`
	UsePrefork      bool   `env:"USE_PREFORK" envDefault:"false"`
	ProfilerEnabeld bool   `env:"PROFILER_ENABLED" envDefault:"false"`
	ProfilerPath    string `env:"PROFILER_PATH" envDefault:"/app/profiling"`
	ProfilerMode    string `env:"PROFILER_MODE" envDefault:"cpu"`
}

func GetConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}
