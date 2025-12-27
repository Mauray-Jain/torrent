package bencode

import (
	"bufio"
	"errors"
	"io"
	"reflect"
	"strconv"
)

// need this if unmarshaling to a map, bcoz we cant directly set on a mapindex
// otherwise can directly pass reflect.Value
type builder struct {
	V   reflect.Value
	Map reflect.Value
	Key reflect.Value
}

func parse(r *bufio.Reader, b *builder) error {
	c, err := r.ReadByte()
	if err != nil {
		return err
	}

	switch {
	case '0' <= c && c <= '9': // string
		err = r.UnreadByte()
		if err != nil {
			return err
		}
		s, err := parseStr(r)
		if err != nil {
			return err
		}
		setString(b.V, s)

	case c == 'i': // int
		i, err := parseInt(r)
		if err != nil {
			return err
		}
		setInt(b.V, i)

	case c == 'l': // list
		v := b.V
		if v.Kind() == reflect.Interface {
			// https://github.com/zeebo/bencode/blob/master/decode.go
			var x []any
			defer func(p reflect.Value) { p.Set(v) }(v) // idk how this works
			v = reflect.ValueOf(&x).Elem()
		}
		if v.Kind() == reflect.Slice && v.IsNil() {
			v.Set(reflect.MakeSlice(v.Type(), 0, 8))
		}
		for i := 0; ; i += 1 {
			c, err := r.ReadByte()
			if err != nil {
				return err
			}
			if c == 'e' {
				break
			}

			if err = r.UnreadByte(); err != nil {
				return err
			}

			val_at_i := getIndex(v, i)
			nb := builder{V: val_at_i}
			if err = parse(r, &nb); err != nil {
				return err
			}
		}

	case c == 'd': // dict
		v := b.V
		isMap := false
		if v.Kind() == reflect.Interface {
			var x map[string]any
			defer func(p reflect.Value) { p.Set(v) }(v)
			v = reflect.ValueOf(&x).Elem()
		}
		if v.Kind() == reflect.Map {
			b.Map = v
			isMap = true
			if v.IsNil() {
				v.Set(reflect.MakeMap(v.Type()))
			}
		}
		for {
			c, err := r.ReadByte()
			if err != nil {
				return err
			}
			if c == 'e' {
				break
			}

			if err = r.UnreadByte(); err != nil {
				return err
			}

			key, err := parseStr(r)
			if err != nil {
				return err
			}

			val_at_key := getKey(v, key)
			if isMap { // if it is a map get an addressable copy of val_at_key
				copy_of_val := reflect.New(val_at_key.Type()).Elem()
				copy_of_val.Set(val_at_key)
				val_at_key = copy_of_val
			}

			nb := builder{V: val_at_key, Map: b.Map, Key: reflect.ValueOf(key)}
			if err = parse(r, &nb); err != nil {
				return err
			}

			if isMap { // after copy of val_at_key is set set it in the map too
				b.Map.SetMapIndex(nb.Key, val_at_key)
			}
		}

	default:
		return errors.New("invalid bencode")
	}

	return nil
}

func parseStr(r *bufio.Reader) (string, error) {
	len_str, err := r.ReadString(':')
	if err != nil {
		return "", err
	}

	length, err := strconv.ParseInt(len_str[:len(len_str)-1], 10, 64)
	if err != nil {
		return "", err
	}
	if length < 0 {
		return "", errors.New("negative length")
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func parseInt(r *bufio.Reader) (int, error) {
	str, err := r.ReadString('e')
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(str[:len(str)-1])
}
