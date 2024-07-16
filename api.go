package phpserialize

import (
	"fmt"

	_ "go4.org/unsafe/assume-no-moving-gc"

	"github.com/trim21/go-phpserialize/internal/encoder"
)

// Marshaler allow users to implement its own encoder.
// **it's return value will not be validated**, please make sure you return valid encoded bytes.
type Marshaler interface {
	MarshalPHP() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalPHP([]byte) error
}

func Marshal(v any) ([]byte, error) {
	return encoder.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	if len(data) == 0 {
		return fmt.Errorf("empty bytes")
	}

	return unmarshal(data, v)
}
