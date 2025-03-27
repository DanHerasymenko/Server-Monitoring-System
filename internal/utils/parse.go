package utils

import (
	"fmt"
	"strconv"
)

func ParseFloat(s string) (float64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string value for float64")
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse float64: %w", err)
	}
	return f, nil
}

func ParseInt64(s string) (int64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string value for int64")
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int64: %w", err)
	}
	return i, nil
}
