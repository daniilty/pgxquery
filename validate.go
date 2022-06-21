package pgxquery

import (
	"errors"
	"reflect"
)

func validateReflectValue(rv reflect.Value) error {
	if !rv.IsValid() {
		return errors.New("invalid value")
	}

	if rv.Kind() != reflect.Struct {
		return errors.New("value must be struct")
	}

	return nil
}
