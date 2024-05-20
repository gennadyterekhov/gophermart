package config

import (
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
	// TODO fix issue-17 https://github.com/gennadyterekhov/gophermart/issues/17
	//addressFlag := flag.String(
	//	"a",
	//	"localhost:8080",
	//	"[address] Net address host:port without protocol",
	//)
	//DBDsnFlag := flag.String(
	//	"d",
	//	"",
	//	"[db dsn] format: `host=%s user=%s password=%s dbname=%s sslmode=%s`",
	//)
	//accrualSystemAddressFlag := flag.String(
	//	"r",
	//	"",
	//	"[accRual_system_address] ",
	//)

	// this breaks tests so i hardcoded defaults
	// flag.Parse()
	//flags := Config{
	//	Addr:       *addressFlag,
	//	DBDsn:      *DBDsnFlag,
	//	AccrualURL: *accrualSystemAddressFlag,
	//}
	const testDBDsn = "host=localhost user=gophermart_user password=gophermart_pass dbname=gophermart_db_test sslmode=disable"

	flags := Config{
		Addr:       "localhost:8080",
		DBDsn:      testDBDsn,
		AccrualURL: "",
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
