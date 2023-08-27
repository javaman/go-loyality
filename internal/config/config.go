package config

import (
	"flag"

	"github.com/caarlos0/env/v9"
)

type Configuration struct {
	Address              string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func Configure() *Configuration {
	config := &Configuration{}

	flag.StringVar(&config.Address, "a", "localhost:8000", "Адрес и порт запуска сервиса")
	flag.StringVar(&config.DatabaseURI, "d", "user=postgres password=iddqd123 host=localhost database=postgres sslmode=disable", "Адрес подключения к базе данных")
	flag.StringVar(&config.AccrualSystemAddress, "r", "localhost:8080", "Адрес системы расчёта начислений")

	env.Parse(config)

	return config
}
