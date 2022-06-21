package pgxquery

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	dbTag = "db"
)

type field struct {
	name        string
	primaryKey  bool
	omitEmpty   bool
	isTableName bool
}

func (f *field) getNamedEqualField() string {
	return f.name + "=" + f.pgxName()
}

func (f *field) pgxName() string {
	return ":" + f.name
}

// waring: this shit panics if rv is not a struct!
func rangeStructFields(rv reflect.Value, f func(reflect.Value, reflect.StructField) bool) {
	l := rv.NumField()
	rt := rv.Type()

	for i := 0; i < l; i++ {
		if !f(rv.Field(i), rt.Field(i)) {
			break
		}
	}
}

func parseField(sf reflect.StructField) (*field, error) {
	const (
		primaryKeyName = "primarykey"
		omitEmptyName  = "omitempty"
	)

	tag, ok := sf.Tag.Lookup(dbTag)
	if !ok {
		return nil, nil
	}

	tags := strings.Split(tag, ",")

	if len(tags) == 0 {
		return nil, fmt.Errorf("%s: struct tag \"db\" must have at least one value", sf.Name)
	}

	isPrimary := false

	i := stringSliceFind(tags, primaryKeyName)
	if i != notFound {
		isPrimary = true

		tags = stringSliceRemove(tags, i)
		if len(tags) == 0 {
			return nil, fmt.Errorf("%s: %s cannot be only struct field tag", sf.Name, primaryKeyName)
		}
	}

	isOmitEmpty := false

	i = stringSliceFind(tags, omitEmptyName)
	if i != notFound {
		isOmitEmpty = true

		tags = stringSliceRemove(tags, i)
		if len(tags) == 0 {
			return nil, fmt.Errorf("%s: %s cannot be only struct field tag", sf.Name, omitEmptyName)
		}
	}

	if len(tags) != 1 {
		return nil, fmt.Errorf("%s: invalid tags len: %d, must be 1(without %s, %s)",
			sf.Name, len(tags), primaryKeyName, omitEmptyName,
		)
	}

	return &field{
		primaryKey:  isPrimary,
		omitEmpty:   isOmitEmpty,
		name:        tags[0],
		isTableName: sf.Type.AssignableTo(reflect.TypeOf(TableName{})),
	}, nil
}
