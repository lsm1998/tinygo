package reflectx

import (
	"reflect"
	"unsafe"
)

// Deprecated: use tinygo.String2Bytes
func String2Bytes(str string) []byte {
	ssh := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&ssh))
}

// Deprecated: use tinygo.Bytes2String
func Bytes2String(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
