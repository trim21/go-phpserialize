package main

import (
	"fmt"
	"reflect"

	"github.com/gookit/goutil/dump"
	"github.com/sirupsen/logrus"
	"github.com/trim21/go-phpserialize"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		ForceQuote:  true,
		// EnvironmentOverrideColors: false,
		DisableTimestamp: false,
		TimestampFormat:  "15:04:05Z07:00",
		FullTimestamp:    false,
		// CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
		// 	return frame.Function, strings.TrimPrefix(file, "C:/Users/Trim21/proj/phpserialize/internal")
		// },
		DisableSorting: false,
		SortingFunc:    nil,
		PadLevelText:   true,
	})

	type Container struct {
		BB     bool     `php:"bb"`
		Value  []string `php:"value"`
		Value2 any      `php:"value2"`
	}

	var c Container
	raw := `a:3:{s:2:"bb";b:1;s:5:"value";a:3:{i:0;s:3:"one";i:1;s:3:"two";i:2;s:1:"q";}s:6:"value2";a:3:{i:0;s:1:"1";i:1;s:1:"2";i:2;s:1:"3";}}`
	err := phpserialize.Unmarshal([]byte(raw), &c)
	if err != nil {
		panic(err)
	}
	expected := []string{"one", "two", "q"}
	if !reflect.DeepEqual(c.Value, expected) {
		dump.P(c.Value)
		dump.P(expected)
		panic("not equal")
	}

	if !c.BB {
		panic("bool parse error")
	}

	fmt.Println("correct")
}
