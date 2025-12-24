package bencode

import (
	"bufio"
	"errors"
	"io"
	"reflect"
	"strconv"
)

func parse(r *bufio.Reader, v reflect.Value) error {
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
		setString(v, s)

	case c == 'i': // int
		i, err := parseInt(r)
		if err != nil {
			return err
		}
		setInt(v, i)

	case c == 'l': // list
		if v.Kind() == reflect.Slice && v.IsNil() {
			v.Set(reflect.MakeSlice(v.Type(), 0, 8))
		}
		i := 0
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

			val_at_i := getIndex(v, i)
			if err = parse(r, val_at_i); err != nil {
				return err
			}
			i += 1
		}

	case c == 'd': // dict
		if v.Kind() == reflect.Map && v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
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
			if err = parse(r, val_at_key); err != nil {
				return err
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

	length, err := strconv.ParseInt(len_str[:len(len_str) - 1], 10, 64)
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
	return strconv.Atoi(str[:len(str) - 1])
}
