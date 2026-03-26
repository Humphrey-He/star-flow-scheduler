package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr string
	MySQLDSN string
}

func Load() Config {
	return Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		MySQLDSN: getEnv("MYSQL_DSN", "root:password@tcp(127.0.0.1:3306)/starflow?parseTime=true&loc=Local"),
	}
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func getEnvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return n
}
