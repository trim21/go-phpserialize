package phpserialize

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/gookit/goutil/dump"
)

const tagName = "php"

func Unmarshal(data []byte, v interface{}) error {
	var t = reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		return errors.New("should be a ptr")
	}

	t = t.Elem()
	if t.Kind() == reflect.Slice {
		return UnmarshalArray(data, v)
	}

	fieldMap := buildFieldsByTagMap(t)

	dump.P(fieldMap)

	/*
		array(
		  'a string value' => 'ff',
		  'a int value' => 31415926,
		  'a bool value' => true,
		  'a float value' => 3.14,
		  662 => 223,
		)
	*/

	/*
	   a:4:{
	      s:14:"a string value";s:2:"ff";
	      s:11:"a int value";i:31415926;
	      s:12:"a bool value";b:1;
	      s:13:"a float value";d:3.14;
	      i:662;i:223;
	   }
	*/

	for offset := 0; offset < len(data); offset++ {
		switch data[offset] {
		case 'o':
			fmt.Println("php object")
		case 'a':
			fmt.Println("php array")
		case 's':
			fmt.Println("php string")
		case 'i':
			fmt.Println("php int")
		case 'b':
			fmt.Println("php bool")
		case 'd':
			fmt.Println("php float")
		}
	}

	return nil
}

func parseASCIIInt(b byte) (int, bool) {
	switch b {
	case '0':
		return 0, true
	case '1':
		return 1, true
	case '2':
		return 2, true
	case '3':
		return 3, true
	case '4':
		return 4, true
	case '5':
		return 5, true
	case '6':
		return 6, true
	case '7':
		return 7, true
	case '8':
		return 8, true
	case '9':
		return 9, true
	}

	return 0, false
}

// consumeString(`14:"a string value";`) => "a string value", 15
func consumeString(p []byte) (s string, offset int, err error) {
	var firstOffset int
	var length int
	for i, b := range p {
		if b == ':' {
			firstOffset = i
			break
		}
		v, ok := parseASCIIInt(b)
		if !ok {
			return "", 0, errors.New("malformed string length prefix")
		}
		length = length*10 + v
	}

	if p[firstOffset+1] != '"' || p[firstOffset+length+2] != '"' {
		return "", 0, errors.New("malformed string")
	}

	s, err = strconv.Unquote(byte2str(p[firstOffset+1 : firstOffset+length+3]))
	offset = firstOffset + length + 4

	return s, offset, err
}

// consumeBool([]byte(`1;`)) => true,2 , nil
// consumeBool([]byte(`0;`)) => true,2 , nil
func consumeBool(p []byte) (v bool, offset int, err error) {
	if p[1] != ';' {
		return false, 0, errors.New("malformed bool")
	}

	switch p[0] {
	case '0':
		return false, 2, nil
	case '1':
		return true, 2, nil
	}

	return false, 0, errors.New("malformed bool")
}

// consumeInt64(`233;`) => int64(233), 15, nil
func consumeInt64(p []byte) (v int64, offset int, err error) {
	var firstOffset int
	for i, b := range p {
		if b == ';' {
			firstOffset = i
			break
		}
		v, ok := parseASCIIInt(b)
		if !ok {
			return 0, 0, errors.New("malformed string length prefix")
		}
		v = v*10 + v
	}

	return v, firstOffset + 1, err
}

//                         Human            json       a1     Head
// map[reflect.Type]map[string]string
var fieldsByTag = sync.Map{}

func buildFieldsByTagMap(rt reflect.Type) map[string]string {
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}

	if v, ok := fieldsByTag.Load(rt); ok {
		return v.(map[string]string)
	}

	m := make(map[string]string, rt.NumField())

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := f.Tag.Get(tagName)
		if v == "" || v == "-" {
			continue
		}
		m[v] = f.Name
	}

	fieldsByTag.Store(rt, m)

	return m
}

func UnmarshalArray(data []byte, v interface{}) error {
	return nil
}
