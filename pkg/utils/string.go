package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ToString(input interface{}) (string, error) {
	switch v := input.(type) {
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v), nil
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", fmt.Errorf("undefined type to convert %T", v)
	}
}

func CheckKeyMatch(pattern, s string) (bool, error) {
	pattern = strings.ReplaceAll(pattern, "/*", "/.*")

	re := regexp.MustCompile(`:[^/]+`)
	pattern = re.ReplaceAllString(pattern, "$1[^/]+$2")

	return regexp.MatchString("^"+pattern+"$", s)
}
