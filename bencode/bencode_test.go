package bencode

import (
	"fmt"
	"strings"
	"testing"
)

type tt struct {
	Test string `bencode:"t"`
	T    string `bencode:"est"`
}

type t2 struct {
	A []string `bencode:"a"`
}

func TestUnmarshal(t *testing.T) {
	s1 := "5:dsd e"
	s2 := "i-834e"
	s3 := "l3:asa4:abcde"
	s4 := "d1:t3:asa3:est2:aae"
	s5 := "d1:al2:bb2:ccee"
	s6 := "d1:ad2:bbi2eee"
	var s string
	var i int
	var l []string
	var l2 [3]string
	var m map[string]string
	var x tt
	var m2 map[string][2]string
	var x2 t2
	var m3 map[string]any
	Unmarshal(strings.NewReader(s1), &s)
	fmt.Println(s)
	Unmarshal(strings.NewReader(s2), &i)
	fmt.Println(i)
	Unmarshal(strings.NewReader(s3), &l)
	fmt.Println(l)
	Unmarshal(strings.NewReader(s3), &l2)
	fmt.Println(l2)
	Unmarshal(strings.NewReader(s4), &x)
	fmt.Println(x)
	Unmarshal(strings.NewReader(s4), &m)
	fmt.Println(m)
	Unmarshal(strings.NewReader(s5), &x2)
	fmt.Println(x2)
	Unmarshal(strings.NewReader(s5), &m2)
	fmt.Println(m2)
	Unmarshal(strings.NewReader(s6), &m3)
	fmt.Println(m3)
}
