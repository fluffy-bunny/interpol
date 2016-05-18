package interpol

import (
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
	case reflect.Struct:
		return structSpawner, nil
	}
	return nil, ErrSpawnerNotFound
}
