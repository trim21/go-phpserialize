package test

import (
	"fmt"
)

type Case struct {
	Name     string
	Data     interface{}
	Expected string `php:"-" json:"-"`
}

func (tc Case) WrappedExpected() string {
	return fmt.Sprintf(`a:2:{s:4:"Name";s:%d:"%s";s:4:"Data";`, len(tc.Name), tc.Name) + tc.Expected + "}"
}

type User struct {
	ID   uint64 `php:"id" json:"id"`
	Name string `php:"name" json:"name"`
}
