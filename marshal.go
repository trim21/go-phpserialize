package phpserialize

import (
	"github.com/trim21/go-phpserialize/internal/encoder"
)

// Marshaler allow users to implement its own encoder.
// **it's return value will not be validated**, please make sure you return valid encoded bytes.
type Marshaler interface {
	MarshalPHP() ([]byte, error)
}

func Marshal(v any) ([]byte, error) {
	return encoder.Marshal(v)
}
