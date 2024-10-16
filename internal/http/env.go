package http

import "os"

func GetEnv(key string, fallback func() string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback()
}
