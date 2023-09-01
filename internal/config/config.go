package config

import (
	"flag"

	"github.com/caarlos0/env/v9"
)

type Configuration struct {
	Address              string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Secret				 string `env:"SECRET"`
}

func Configure() *Configuration {
	config := &Configuration{}

	flag.StringVar(&config.Address, "a", "localhost:8000", "Адрес и порт запуска сервиса")
	flag.StringVar(&config.DatabaseURI, "d", "user=postgres password=iddqd123 host=localhost database=postgres sslmode=disable", "Адрес подключения к базе данных")
	flag.StringVar(&config.AccrualSystemAddress, "r", "http://localhost:8080", "Адрес системы расчёта начислений")
	flag.StringVar(&config.Secret, "s", "iddqd", "Секрет для JWT")

	env.Parse(config)

	return config
}
