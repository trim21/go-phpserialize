package phpserialize

import (
	_ "go4.org/unsafe/assume-no-moving-gc"

	"github.com/trim21/go-phpserialize/internal/encoder"
)

func Marshal(v any) ([]byte, error) {
	return encoder.Marshal(v)
}
