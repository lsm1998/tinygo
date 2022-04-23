package reflectx

import (
	"reflect"
	"unsafe"
)

func String2Bytes(str string) []byte {
	ssh := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&ssh))
}

func Bytes2String(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
