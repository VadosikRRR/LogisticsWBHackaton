package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetString(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("invalid int env %s=%q, using fallback %d", key, val, fallback)
		return fallback
	}
	return parsed
}

func GetFloat(key string, fallback float64) float64 {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Printf("invalid float env %s=%q, using fallback %.3f", key, val, fallback)
		return fallback
	}
	return parsed
}

func GetBool(key string, fallback bool) bool {
	val := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if val == "" {
		return fallback
	}
	switch val {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		log.Printf("invalid bool env %s=%q, using fallback %t", key, val, fallback)
		return fallback
	}
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("invalid duration env %s=%q, using fallback %s", key, val, fallback)
		return fallback
	}
	return parsed
}
