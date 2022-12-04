package helpers

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func getEnvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}

func GetEnvInt(key string) int {
	s, err := getEnvStr(key)
	if err != nil {
		return 0
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func GetEnvBool(key string) bool {
	s, err := getEnvStr(key)
	if err != nil {
		return false
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return v
}

func GetEnvTimeDuration(key string) time.Duration {
	s, err := getEnvStr(key)
	if err != nil {
		return 0
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return v
}

func GetEnvListString(key string) []string {
	s, err := getEnvStr(key)
	if err != nil {
		return nil
	}
	return strings.Split(s, ",")
}
