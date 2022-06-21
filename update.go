package pgxquery

import (
	"errors"
	"reflect"
	"strings"
)

// GenerateNamedUpdate - generate basic update query from given value.
// value must be struct or pointer to a struct.
// Or else it will panic.
func GenerateNamedUpdate(v interface{}) (string, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return "", errors.New("nil pointer value")
		}

		rv = rv.Elem()
	}

	err := validateReflectValue(rv)
	if err != nil {
		return "", err
	}

	return generateUpdate(rv)
}

func generateUpdate(rv reflect.Value) (string, error) {
	const separator = ", "

	var (
		err       error
		tableName string

		args      = []string{}
		primaries = []string{}
	)

	rangeStructFields(rv, func(v reflect.Value, sf reflect.StructField) bool {
		var field *field

		field, err = parseField(sf)
		if err != nil {
			return false
		}

		if field.primaryKey {
			primaries = append(primaries, field.getNamedEqualField())

			return true
		}

		if field.isTableName {
			tableName = field.name

			return true
		}

		// omit zero value field
		if field.omitEmpty {
			if v.IsZero() {
				return true
			}
		}

		args = append(args, field.getNamedEqualField())

		return true
	})

	if err != nil {
		return "", err
	}

	if len(primaries) == 0 {
		return "", errors.New("no primary key")
	}

	if tableName == "" {
		return "", errors.New("add TableName to your model")
	}

	if len(args) == 0 {
		return "", ErrEmptyModel
	}

	argsStr := strings.Join(args, separator)
	primariesStr := strings.Join(primaries, " and ")

	return "update " + tableName + " set " + argsStr + " where " + primariesStr + ";", nil
}
