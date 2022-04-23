package reflectx

import (
	"fmt"
	"testing"
)

func TestBytes2String(t *testing.T) {
	str := Bytes2String([]byte("hello"))
	fmt.Println(str)
}

func TestString2Bytes(t *testing.T) {
	str := "hello"
	bytes := String2Bytes(str)
	fmt.Println(bytes)
}
