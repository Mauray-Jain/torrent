package bencode

import (
	"fmt"
	"reflect"
	"strings"
)

func setString(v reflect.Value, s string) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Interface:
		v.Set(reflect.ValueOf(s))
	default:
		panic("not a string")
	}
}

func setInt(v reflect.Value, i int) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(i))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(uint64(i))
	case reflect.Interface:
		v.Set(reflect.ValueOf(i))
	default:
		panic("not an int")
	}
}

func getIndex(v reflect.Value, i int) reflect.Value {
	switch v.Kind() {
	case reflect.Array:
		return v.Index(i)
	case reflect.Slice:
		v_cap := v.Cap()
		if i >= v_cap {
			if v_cap < 8 {
				v_cap = 8
			} else {
				v_cap *= 2
			}
			new_slice := reflect.MakeSlice(v.Type(), v.Len(), v_cap)
			_ = reflect.Copy(new_slice, v)
			v.Set(new_slice)
		}
		if i >= v.Len() {
			v.SetLen(i + 1)
		}
		return v.Index(i)
	default:
		panic("not an array or slice")
	}
}

func getKey(v reflect.Value, s string) reflect.Value {
	switch v.Kind() {
	case reflect.Map:
		panic("TODO")
		// key := reflect.ValueOf(s)
		// elem := v.MapIndex(key)
		// if !elem.IsValid() {
		// 	v.SetMapIndex(key, reflect.Zero(v.Type().Elem()))
		// 	elem = v.MapIndex(key)
		// }
		// fmt.Println(s, elem.CanSet())
		// return elem

	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			tf := t.Field(i)
			tagname, ok := tf.Tag.Lookup("bencode")
			if !ok {
				continue
			}
			if tagname == s || strings.ToLower(tf.Name) == s {
				return v.Field(i)
			}
		}
		panic("not in struct field")

	default:
		panic("not a map or struct")
	}
}
