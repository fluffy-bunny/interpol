package interpol

import (
	"fmt"
	"strconv"
)

func getStringFromInterface(v interface{}) ([]byte, error) {
	// Check some basic types
	switch v := v.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	case int:
		return []byte(strconv.Itoa(v)), nil
	case int32:
		return []byte(strconv.Itoa(int(v))), nil
	case int64:
		return []byte(strconv.FormatInt(v, 10)), nil
	case uint:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case uint32:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case uint64:
		return []byte(strconv.FormatUint(v, 10)), nil
	}

	// check if an item implements the fmt.Stringer
	if str, ok := v.(fmt.Stringer); ok {
		return []byte(str.String()), nil
	}

	// last resort
	return []byte(fmt.Sprintf("%v", v)), nil
}
