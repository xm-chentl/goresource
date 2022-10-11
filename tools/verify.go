package tools

import "reflect"

// IsEmpty support: string、int 8,16,32,64
func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		vv := v.(string)
		return vv == ""
	case reflect.Int:
		vv := v.(int)
		return vv == 0
	case reflect.Int8:
		vv := v.(int8)
		return vv == 0
	case reflect.Int16:
		vv := v.(int16)
		return vv == 0
	case reflect.Int32:
		vv := v.(int32)
		return vv == 0
	case reflect.Int64:
		vv := v.(int64)
		return vv == 0
	}

	return false
}
