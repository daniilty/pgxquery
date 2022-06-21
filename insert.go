package pgxquery

import (
	"errors"
	"reflect"
	"strings"
)

// GenerateNamedInsert - generate basic insert query from given value.
// value must be struct or pointer to a struct.
// Or else it will panic.
func GenerateNamedInsert(v interface{}) (string, error) {
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

	return generateAdd(rv)
}

func generateAdd(rv reflect.Value) (string, error) {
	const separator = ", "

	var (
		err       error
		tableName string

		fieldNames = []string{}
		values     = []string{}
	)

	rangeStructFields(rv, func(_ reflect.Value, sf reflect.StructField) bool {
		var field *field

		field, err = parseField(sf)
		if err != nil {
			return false
		}

		if field == nil {
			return true
		}

		if field.isTableName {
			tableName = field.name

			return true
		}

		fieldNames = append(fieldNames, field.name)
		values = append(values, field.pgxName())

		return true
	})

	if err != nil {
		return "", err
	}

	if tableName == "" {
		return "", errors.New("add TableName to your model")
	}

	if len(fieldNames) == 0 {
		return "", ErrEmptyModel
	}

	fieldNamesStr := strings.Join(fieldNames, separator)
	valuesStr := strings.Join(values, separator)

	return "insert into " + tableName + inBrackets(fieldNamesStr) + " values" + inBrackets(valuesStr) + ";", nil
}
