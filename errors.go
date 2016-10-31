package framework

import (
	"errors"
	"fmt"
)

var errMissing = errors.New("key is missing")

type errNotAString struct {
	value interface{}
}

func (e errNotAString) Error() string {
	return fmt.Sprintf("value %v is not a string", e.value)
}
