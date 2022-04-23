package reflectx

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// Struct2Map 结构体转为map[string]string
func Struct2Map(val interface{}) map[string]string {
	typeOf := reflect.TypeOf(val)
	valueOf := reflect.ValueOf(val)
	for typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	count := typeOf.NumField()
	result := make(map[string]string, count)
	for i := 0; i < count; i++ {
		putFieldKeyValue(result, typeOf.Field(i), valueOf.Field(i))
	}
	return result
}

// Struct2UrlValues 结构体转为url.Values
func Struct2UrlValues(val interface{}) url.Values {
	struct2Map := Struct2Map(val)
	result := make(url.Values, len(struct2Map))
	for k, v := range struct2Map {
		result[k] = []string{v}
	}
	return result
}

func putFieldKeyValue(m map[string]string, field reflect.StructField, value reflect.Value) {
	tags := strings.Split(field.Tag.Get("json"), ",")
	if len(tags) == 0 || tags[0] == "-" {
		return
	}
	strVal, ok := getStringValue(value)
	for i := 1; i < len(tags); i++ {
		if ok && tags[i] == "omitempty" {
			return
		}
	}
	m[tags[0]] = strVal
}

func getStringValue(value reflect.Value) (string, bool) {
	switch value.Kind() {
	case reflect.String:
		v := value.String()
		return v, v == ""
	case reflect.Int:
		v := value.Int()
		return fmt.Sprintf("%d", v), v == 0
	case reflect.Int8:
		v := value.Int()
		return fmt.Sprintf("%d", v), v == 0
	case reflect.Int16:
		v := value.Int()
		return fmt.Sprintf("%d", v), v == 0
	case reflect.Int32:
		v := value.Int()
		return fmt.Sprintf("%d", v), v == 0
	case reflect.Int64:
		v := value.Int()
		return fmt.Sprintf("%d", v), v == 0
	case reflect.Uint:
		v := value.Int()
		return fmt.Sprintf("%d", v), v == 0
	case reflect.Uint8:
		v := value.Int()
		return fmt.Sprintf("%d", v), v == 0
	case reflect.Uint16:
		v := value.Int()
		return fmt.Sprintf("%d", v), v == 0
	case reflect.Uint32:
		v := value.Int()
		return fmt.Sprintf("%d", v), v == 0
	case reflect.Uint64:
		v := value.Int()
		return fmt.Sprintf("%d", v), v == 0
	case reflect.Float32:
		v := value.Float()
		return fmt.Sprintf("%f", v), v == 0
	case reflect.Float64:
		v := value.Float()
		return fmt.Sprintf("%f", v), v == 0
	case reflect.Bool:
		v := value.Bool()
		return fmt.Sprintf("%v", v), !v
	case reflect.Chan:
	case reflect.Struct:
	case reflect.Uintptr:
	case reflect.Map:
		values := make([]string, 0, value.Len())
		iter := value.MapRange()
		for iter.Next() {
			k, _ := getStringValue(iter.Key())
			v, _ := getStringValue(iter.Value())
			values = append(values, k+":"+v)
		}
		return "{" + strings.Join(values, ",") + "}", value.Len() == 0
	case reflect.Array:
		arr := make([]string, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			arrVal, _ := getStringValue(value.Index(i))
			arr = append(arr, arrVal)
		}
		return strings.Join(arr, ","), value.Len() == 0
	case reflect.Slice:
		arr := make([]string, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			arrVal, _ := getStringValue(value.Index(i))
			arr = append(arr, arrVal)
		}
		return strings.Join(arr, ","), value.Len() == 0
	}
	return "", true
}
