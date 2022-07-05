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

	// logrus.SetReportCaller(true)

	type Container struct {
		Value map[string]int `php:"value"`
	}

	var c Container
	raw := `a:1:{s:5:"value";a:2:{s:3:"one";i:1;s:3:"two";i:2;}}`
	err := phpserialize.Unmarshal([]byte(raw), &c)
	if err != nil {
		panic(err)
	}
	expected := map[string]int{"one": 1, "two": 2}
	if !reflect.DeepEqual(c.Value, expected) {
		dump.P(c.Value)
		dump.P(expected)
		panic("not equal")
	}

	fmt.Println("correct")
}
