package iphotos

import (
	"encoding/json"
	"fmt"
)

type Int64T interface {
	uint |
		uint8 |
		uint16 |
		uint32 |
		uint64 |
		int |
		int8 |
		int16 |
		int32 |
		int64 |
		float32 |
		float64
}

func uintToInt64[T Int64T](vs ...T) int64 {
	if len(vs) == 0 {
		return 0
	}
	return int64(vs[0])
}

func anyToInt64(vs any) (int64, bool) {
	var result int64
	switch v := vs.(type) {
	case []uint:
		result = uintToInt64(v...)
	case []uint8:
		result = uintToInt64(v...)
	case []uint16:
		result = uintToInt64(v...)
	case []uint32:
		result = uintToInt64(v...)
	case []uint64:
		result = uintToInt64(v...)
	case []int8:
		result = uintToInt64(v...)
	case []int16:
		result = uintToInt64(v...)
	case []int32:
		result = uintToInt64(v...)
	case []int64:
		result = uintToInt64(v...)
	case []float32:
		result = uintToInt64(v...)
	case []float64:
		result = uintToInt64(v...)
	case uint:
		result = uintToInt64(v)
	case uint8:
		result = uintToInt64(v)
	case uint16:
		result = uintToInt64(v)
	case uint32:
		result = uintToInt64(v)
	case uint64:
		result = uintToInt64(v)
	case int8:
		result = uintToInt64(v)
	case int16:
		result = uintToInt64(v)
	case int32:
		result = uintToInt64(v)
	case int64:
		result = uintToInt64(v)
	case float32:
		result = uintToInt64(v)
	case float64:
		result = uintToInt64(v)
	default:
		return 0, false
	}
	return result, result > 0
}
func anyToInt64Value(vs any) int64 {
	v, _ := anyToInt64(vs)
	return v
}

func anyToStrings(vs any) []string {
	switch v := vs.(type) {
	case []string:
		return v
	case string:
		return []string{v}
	case []any:
		if by, err := json.Marshal(v); err == nil {
			var result []string
			json.Unmarshal(by, &result)
			return result
		}
	case any:
		return []string{
			fmt.Sprintf("%v", v),
		}
	}
	return []string{}
}
