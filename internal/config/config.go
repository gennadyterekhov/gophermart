package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr       string
	IsGzip     bool
	DBDsn      string
	AccrualURL string
}

var ServerConfig *Config = getConfig()

func getConfig() *Config {
	addressFlag := flag.String(
		"a",
		"localhost:8080",
		"[address] Net address host:port without protocol",
	)
	DBDsnFlag := flag.String(
		"d",
		"",
		"[db dsn] format: `host=%s user=%s password=%s dbname=%s sslmode=%s`",
	)
	accrualSystemAddressFlag := flag.String(
		"r",
		"",
		"[accRual_system_address] ",
	)

	flag.Parse()

	flags := Config{
		Addr:       *addressFlag,
		DBDsn:      *DBDsnFlag,
		AccrualURL: *accrualSystemAddressFlag,
	}

	overwriteWithEnv(&flags)

	return &flags
}

func overwriteWithEnv(flags *Config) {
	flags.Addr = getAddress(flags.Addr)
	flags.DBDsn = getDBDsn(flags.DBDsn)
	flags.AccrualURL = getAccrualURL(flags.AccrualURL)
}

func getAddress(current string) string {
	return getStringFromEnvOrFallback("RUN_ADDRESS", current)
}

func getDBDsn(current string) string {
	return getStringFromEnvOrFallback("DATABASE_URI", current)
}

func getAccrualURL(current string) string {
	return getStringFromEnvOrFallback("ACCRUAL_SYSTEM_ADDRESS", current)
}

func getStringFromEnvOrFallback(envKey string, fallback string) string {
	fromEnv, ok := os.LookupEnv(envKey)
	if ok {
		return fromEnv
	}

	return fallback
}
