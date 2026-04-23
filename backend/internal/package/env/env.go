package env

import (
	"os"
	"strconv"
)

func GetWithDefault[T any](key string, defaultValue T) T {
	valStr, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	var result any

	switch any(defaultValue).(type) {
	case string:
		result = valStr
	case int:
		parsed, err := strconv.Atoi(valStr)
		if err != nil {
			return defaultValue
		}
		result = parsed
	case bool:
		parsed, err := strconv.ParseBool(valStr)
		if err != nil {
			return defaultValue
		}
		result = parsed
	default:
		return defaultValue
	}

	return result.(T)
}
