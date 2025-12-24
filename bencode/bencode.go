// Heavily borrowed and inspired from https://github.com/jackpal/bencode-go
package bencode

import (
	"bufio"
	"errors"
	"io"
	"reflect"
)

func Unmarshal(r io.Reader, v any) error {
	val := reflect.ValueOf(v)
	if v == nil || val.Kind() != reflect.Pointer {
		return errors.New("v is not a pointer")
	}
	reader := bufio.NewReader(r)
	return parse(reader, val.Elem())
}
