//go:build !1.23rc2

package phpserialize

import (
	"fmt"

	"github.com/trim21/go-phpserialize/internal/decoder"
	"github.com/trim21/go-phpserialize/internal/encoder"
)

// Marshaler allow users to implement its own encoder
// **it's return value will not be validated**, please make sure you return valid encoded bytes.
type Marshaler = encoder.Marshaler

type Unmarshaler = decoder.Unmarshaler

func Marshal(v any) ([]byte, error) {
	return encoder.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	if len(data) == 0 {
		return fmt.Errorf("empty bytes")
	}

	return unmarshal(data, v)
}
