package utils

import (
	"os"
	"strconv"
	"strings"
)

type InstagramConfig struct {
	Username    string
	Password    string
	Donors      []string
	SubsPerHour int64
	SubInterval int64
	Condition_1 bool
	Condition_2 bool
}

type Config struct {
	Instagram InstagramConfig
	DebugMode bool
}

func NewConfig() *Config {
	return &Config{
		Instagram: InstagramConfig{
			Username:    getEnv("INSTAGRAM_USERNAME", ""),
			Password:    getEnv("INSTAGRAM_PASSWORD", ""),
			Donors:      getEnvAsStringSlice("INSTAGRAM_DONORS", make([]string, 0)),
			SubsPerHour: getEnvAsInt("INSTAGRAM_SUBS_PER_HOUR", 2),
			SubInterval: getEnvAsInt("INSTAGRAM_INTERVAL", 100),
			Condition_1: getEnvAsBool("INSTAGRAM_CONDITION_1", false),
			Condition_2: getEnvAsBool("INSTAGRAM_CONDITION_2", false),
		},
		DebugMode: getEnvAsBool("DEBUGMODE", true),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valString := getEnv(key, "")
	if value, err := strconv.ParseBool(valString); err == nil {
		return value
	}

	return defaultValue
}

func getEnvAsInt(key string, defaultValue int64) int64 {
	valString := getEnv(key, "")
	if value, err := strconv.ParseInt(valString, 10, 32); err == nil {
		return value
	}

	return defaultValue
}

func getEnvAsStringSlice(key string, defaultValue []string) []string {
	valString := getEnv(key, "")
	return strings.Split(valString, ",")
}
