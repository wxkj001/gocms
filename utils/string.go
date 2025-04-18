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
