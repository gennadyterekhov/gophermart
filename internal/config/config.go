package config

import (
	"flag"
	"os"

	"github.com/gennadyterekhov/gophermart/internal/logger"
)

type Config struct {
	Addr       string
	IsGzip     bool
	DBDsn      string
	AccrualURL string
}

var ServerConfig *Config = &Config{
	Addr:       "localhost:8888",
	DBDsn:      "host=localhost user=gophermart_user password=gophermart_pass dbname=gophermart_db_test sslmode=disable",
	AccrualURL: "http://localhost:8080",
}

func Init() {
	ServerConfig = getConfig()
}

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
	logger.CustomLogger.Debugln("flags before envs", flags.Addr, flags.DBDsn, flags.AccrualURL)
	overwriteWithEnv(&flags)
	logger.CustomLogger.Debugln("flags after envs", flags.Addr, flags.DBDsn, flags.AccrualURL)

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
