package interpol

// func used to retrieve the value from an object (struct or map)
type getterFunc func(string) ([]byte, error)
type getterFuncSpawner func(interface{}) (getterFunc, error)

func getterSpawnerSelector(v interface{}) (getterFuncSpawner, error) {
	switch v.(type) {
	case map[string]string:
		return mapStringStringSpawner, nil
	case map[string][]byte:
		return mapStringByteSpawner, nil
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
