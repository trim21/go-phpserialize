package phpserialize

import (
	"github.com/trim21/go-phpserialize/internal/encoder"
)

func Marshal(v any) ([]byte, error) {
	return encoder.Marshal(v)
}
