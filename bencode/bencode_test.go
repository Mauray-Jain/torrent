package bencode

import (
	"fmt"
	"strings"
	"testing"
)

type tt struct {
	Test string `bencode:"t"`
	T string `bencode:"est"`
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
	var s string
	var i int
	var l []string
	var l2 [3]string
	var m map[string]string
	// var m2 map[string][2]string
	var x tt
	var x2 t2
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
	// Unmarshal(strings.NewReader(s5), &m2) // fails
	// fmt.Println(m2)
}
