package main

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Port int
}

func GetEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func LoadEnvs(c *Config) {
	flag.IntVar(&c.Port, "port", GetEnvInt("PORT", 4221), "Port to listen on")
}
