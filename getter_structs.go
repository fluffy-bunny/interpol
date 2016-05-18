package interpol

import (
	"reflect"
	"strings"
)

func structSpawner(v interface{}) (getterFunc, error) {
	refVal := reflect.ValueOf(v)
	return func(key string) ([]byte, error) {
		value, err := structWalk(key, refVal)
		if err != nil {
			return nil, err
		}
		return getStringFromInterface(value.Interface())
	}, nil
}

// structWalk walks the struct recursively, splitting the key
// at dots and moving further down the tree.
func structWalk(key string, v reflect.Value) (reflect.Value, error) {
	// Check if it's a struct first
	if v.Kind() != reflect.Struct {
		return v, ErrStructKeyNotFound
	}

	keys := strings.SplitN(key, ".", 2)

	field := v.FieldByName(keys[0])
	if !field.IsValid() {
		return field, ErrStructKeyNotFound
	}
	if keys[1] == `` {
		return field, nil
	}

	return structWalk(keys[1], field)
}
