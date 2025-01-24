package go3s

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func toParams(s any) (params url.Values) {
	params = url.Values{}
	typ := reflect.TypeOf(s)
	val := reflect.ValueOf(s)
	kind := val.Kind()
	// s may be nil, invalid value
	if !val.IsValid() {
		return
	}
	if kind != reflect.Struct {
		panic(fmt.Sprintf("solscan: input kind %v is not a struct", kind))
	}
	// if kind == reflect.Pointer || kind == reflect.UnsafePointer {
	// 	return toParams(val.Interface())
	// }
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		name := field.Name
		v := val.FieldByName(name)
		if !v.IsValid() {
			panic(fmt.Sprintf("solscan: field %v is invalid", name))
		}
		tag, ok := field.Tag.Lookup("json")
		if ok {
			// field tag name may be omitempty, so can not use strings.Contains(tag, "omitempty").
			tags := strings.Split(tag, ",")
			if len(tags) > 1 &&
				strings.Contains(tags[1], "omitempty") && v.IsZero() {
				continue
			}
			name = strings.Trim(tags[0], " ")
		}

		tag, ok = field.Tag.Lookup("default")
		if ok && v.IsZero() {
			params.Add(name, tag)
			continue
		}

		switch v.Kind() {
		case reflect.String, reflect.Bool, reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr, reflect.Float64, reflect.Float32:
			params.Add(name, fmt.Sprintf("%v", v.Interface()))
		case reflect.Slice, reflect.Array:
			l := v.Len()
			for i := 0; i < l; i++ {
				params.Add(name+"[]", fmt.Sprintf("%v", v.Index(i).Interface()))
			}
		case reflect.Invalid:
			panic(fmt.Sprintf("solscan: field %v is invalid", name))
		default: // reflect.Interface
			panic(fmt.Sprintf("solscan: field %v is interface", name))
		}
	}
	return
}
