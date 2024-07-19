package common

import (
	"syscall"
)

//takes key and fallback string and returns the key if get env return true
func EnvString(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}

	return fallback
}