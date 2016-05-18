package interpol

import (
	"errors"
	"reflect"
)

// func used to retrieve the value from an object (struct or map)
type getterFunc func(string) ([]byte, error)
type getterFuncSpawner func(interface{}) (getterFunc, error)

func getterSpawnerSelector(v interface{}) (getterFuncSpawner, error) {
	val := reflect.ValueOf(v)
	switch val.Type().Kind() {
	case reflect.Map:
		// Map type, get selector based on the value
		return getMapSelector(val.Type())
	}
	return nil, ErrSpawnerNotFound
}

// getMapSelector retrieves the spawner based on the map's value type.
func getMapSelector(v reflect.Type) (getterFuncSpawner, error) {
	if v.Key().Kind() != reflect.String {
		return nil, errors.New("Non-string keyed maps are not supported!")
	}
	switch v.Elem().Kind() {
	case reflect.String:
		// simple string-string map
		return mapStringStringSpawner, nil
	case reflect.Slice:
		// map of slice
		switch v.Elem().Elem().Kind() {
		// safe to assume that's []byte
		case reflect.Uint8:
			return mapStringByteSpawner, nil
		}
	}
	return nil, ErrSpawnerNotFound
}

// mapStringString is used for map[string]string lookup.
func mapStringStringSpawner(v interface{}) (getterFunc, error) {
	m := v.(map[string]string)
	return func(key string) ([]byte, error) {
		value, ok := m[key]
		if !ok {
			return nil, ErrMapKeyNotFound
		}
		return []byte(value), nil
	}, nil
}

// mapStringByte is used for map[string][]byte lookups.
func mapStringByteSpawner(v interface{}) (getterFunc, error) {
	m := v.(map[string][]byte)
	return func(key string) ([]byte, error) {
		value, ok := m[key]
		if !ok {
			return nil, ErrMapKeyNotFound
		}
		return value, nil
	}, nil
}
