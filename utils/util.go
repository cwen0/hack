package utils

import (
	"github.com/juju/errors"
	"strings"
)

// MatchInArray returns true if the given string value is in the array.
func MatchInArray(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

// GetIP returns ip from addr
func GetIP(addr string) (string, error) {
	arrs := strings.Split(addr, ":")
	if len(arrs) < 2 {
		return "", errors.NotValidf("%s", addr)
	}

	return arrs[0], nil
}