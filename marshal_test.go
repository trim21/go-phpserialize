package phpserialize_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trim21/go-phpserialize"
)

type Inner struct {
	V int    `php:"v"`
	S string `php:"a long string name replace field name"`
}

type TestData struct {
	Users []User                  `php:"users"`
	Obj   Inner                   `php:"obj"`
	B     bool                    `php:"ok"`
	Map   map[int]struct{ V int } `php:"map"`
}

type User struct {
	ID   uint64 `php:"id"`
	Name string `php:"name"`
}

func TestMarshal(t *testing.T) {
	t.Skip()
	fmt.Println("TestMarshal start")
	var data = TestData{
		Users: []User{
			{ID: 1, Name: "sai"},
			{ID: 2, Name: "trim21"},
		},
		B:   false,
		Obj: Inner{V: 2, S: "vvv"},

		Map: map[int]struct{ V int }{7: {V: 4}},
	}

	_, err := phpserialize.Marshal(data)
	require.NoError(t, err)
	fmt.Println("TestMarshal end")
}

//
