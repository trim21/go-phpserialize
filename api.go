//go:build !go1.23

package phpserialize

import (
	"fmt"

	"github.com/trim21/go-phpserialize/internal/encoder"
)

func Marshal(v any) ([]byte, error) {
	return encoder.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	if len(data) == 0 {
		return fmt.Errorf("empty bytes")
	}

	return unmarshal(data, v)
}

type Unmarshaler interface {
	UnmarshalPHP([]byte) error
}
