package reflectx

import (
	"reflect"
	"unsafe"
)

// SizeOf 获取一个对象实际占用字节大小
func SizeOf(obj interface{}) int64 {
	return sizeof(reflect.ValueOf(obj))
}

func sizeof(v reflect.Value) int64 {
	var sum int64
	switch v.Kind() {
	case reflect.Map:
		keys := v.MapKeys()
		for i := 0; i < len(keys); i++ {
			key := keys[i]
			sum += sizeof(key)
			sum += sizeof(v.MapIndex(key))
		}
		return sum
	case reflect.Slice, reflect.Array:
		for i, n := 0, v.Len(); i < n; i++ {
			sum += sizeof(v.Index(i))
		}
		return sum
	case reflect.String:
		return int64(len(v.String()))
	case reflect.Ptr, reflect.Interface:
		p := (*[]byte)(unsafe.Pointer(v.Pointer()))
		if p == nil {
			return 0
		}
		return sizeof(v.Elem())
	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			sum += sizeof(v.Field(i))
		}
		return sum
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Int:
		return int64(int(v.Type().Size()))
	default:
		return 0
	}
}
