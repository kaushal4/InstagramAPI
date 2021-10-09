package utility

import (
	"errors"
	"strconv"
	"strings"
)

func GetIDFromString(id string) (int, error) {
	if i, err := strconv.Atoi(id); err == nil {
		return i, nil
	} else {
		return 0, errors.New("invalid ID")
	}
}

func ExtractID(path, base string) (int, error) {
	id := strings.TrimPrefix(path, base)
	return GetIDFromString(id)
}
