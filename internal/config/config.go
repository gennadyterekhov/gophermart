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
	gzipFlag := flag.Bool(
		"g",
		true,
		"[gzip] use gzip",
	)
	DBDsnFlag := flag.String(
		"d",
		"",
		"[db dsn] format: `host=%s user=%s password=%s dbname=%s sslmode=%s` if empty, ram is used",
	)

	accrualSystemAddressFlag := flag.String(
		"r",
		"",
		"[accRual_system_address] ",
	)

	flag.Parse()

	flags := Config{
		Addr:       *addressFlag,
		IsGzip:     *gzipFlag,
		DBDsn:      *DBDsnFlag,
		AccrualURL: *accrualSystemAddressFlag,
	}

	overwriteWithEnv(&flags)

	return &flags
}

func overwriteWithEnv(flags *Config) {
	flags.Addr = getAddress(flags.Addr)
	flags.IsGzip = isGzip(flags.IsGzip)
	flags.DBDsn = getDBDsn(flags.DBDsn)
	flags.AccrualURL = getAccrualURL(flags.AccrualURL)
}

func getAddress(current string) string {
	rawAddress, ok := os.LookupEnv("RUN_ADDRESS")
	if ok {
		return rawAddress
	}

	return current
}

func isGzip(gzip bool) bool {
	fromEnv, ok := os.LookupEnv("GZIP")
	if ok && (fromEnv == "true" || fromEnv == "TRUE" || fromEnv == "True" || fromEnv == "1") {
		return true
	}
	if ok {
		return false
	}

	return gzip
}

func getDBDsn(current string) string {
	raw, ok := os.LookupEnv("DATABASE_URI")
	if ok {
		return raw
	}

	return current
}

func getAccrualURL(current string) string {
	raw, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if ok {
		return raw
	}

	return current
}
