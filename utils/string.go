package utils

import "strconv"

func ToString(d any) string {
	switch d.(type) {
	case string:
		return d.(string)
	case int:
		return strconv.Itoa(d.(int))
	case int64:
		return strconv.FormatInt(d.(int64), 10)
	case float64:
		return strconv.FormatFloat(d.(float64), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(d.(float32)), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(d.(bool))
	default:
		return ""
	}
}
func ToInt(d any) int {
	switch d.(type) {
	case string:
		i, _ := strconv.Atoi(d.(string))
		return i
	case int:
		return d.(int)
	case int64:
		return int(d.(int64))
	case uint:
		return int(d.(uint))
	case uint64:
		return int(d.(uint64))
	case uint32:
		return int(d.(uint32))
	case float64:
		return int(d.(float64))
	case float32:
		return int(d.(float32))
	case bool:
		if d.(bool) {
			return 1
		} else {
			return 0
		}
	default:
		return 0
	}
}
