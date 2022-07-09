package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	j, err := json.Marshal(MapOnly{
		Map: map[string]int64{"one": 1},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))

}

type Item struct {
	V int `json:"v" php:"v"`
}

// map in struct is a direct ptr
type MapOnly struct {
	Map map[string]int64 `php:"map" json:"map"`
}
