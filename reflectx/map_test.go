package reflectx

import (
	"fmt"
	"testing"
)

type Student struct {
	Name   string  `json:"name"`
	Age    int     `json:"age,omitempty"`
	Adders []int32 `json:"adders"`

	Names [2]string `json:"names"`

	F1 float32 `json:"f_1"`
	F2 float64 `json:"f_2"`
	F3 int8    `json:"-"`

	Bl bool `json:"bl"`

	Scores map[string]int32 `json:"scores"`
}

func TestStruct2Map(t *testing.T) {
	s := Student{Name: "lsm"}
	s.F1 = 1.1
	s.F2 = 1.2
	s.F3 = 1
	s.Age = 100
	s.Adders = []int32{1, 2, 4}
	s.Names = [2]string{"hello", "master"}
	s.Scores = map[string]int32{
		"语文": 90,
		"数学": 90,
	}
	fmt.Println(Struct2Map(&s))
	// fmt.Println(Struct2UrlValues(&s))
}

func BenchmarkStruct2Map(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := Student{Name: "lsm"}
		s.F1 = 1.1
		s.F2 = 1.2
		s.F3 = 1
		s.Adders = []int32{1, 2, 4}
		Struct2Map(&s)
	}
}
