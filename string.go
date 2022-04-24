package tinygo

import (
	"github.com/google/uuid"
	"reflect"
	"strings"
	"unsafe"
)

func UUID() string {
	return uuid.New().String()
}

func IsBlank(str string) bool {
	return len(strings.ReplaceAll(str, " ", "")) == 0
}

func String2Bytes(str string) []byte {
	ssh := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&ssh))
}

func Bytes2String(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
