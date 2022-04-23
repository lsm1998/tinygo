package reflectx

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestSizeOf(t *testing.T) {
	var a int32 = 100
	fmt.Println(SizeOf(a))

	var b int64 = 100
	fmt.Println(SizeOf(b))

	fmt.Println(unsafe.Sizeof(a))
	fmt.Println(unsafe.Sizeof(b))

	type Demo1 struct {
		Name string `json:"name"` // 16
		Age  uint32 `json:"age"`  // 4
	}
	type Demo2 struct {
		D1   Demo1 `json:"d_1"` //
		Time int64 `json:"time"`
	}
	fmt.Println(unsafe.Sizeof(new(Demo2)))
	fmt.Println(SizeOf("hello"))
	fmt.Println(SizeOf(new(Demo1)))
	fmt.Println(SizeOf(new(Demo2)))
}
