package utility

import (
	"errors"
	"strconv"
)

func GetIDFromString(id string) (int, error) {
	if i, err := strconv.Atoi(id); err == nil {
		return i, nil
	} else {
		return 0, errors.New("invalid ID")
	}
}
