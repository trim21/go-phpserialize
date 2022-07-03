package phpserialize

import "github.com/trim21/go-phpserialize/internal/decoder"

func Unmarshal(data []byte, v any) error {
	return decoder.Unmarshal(data, v)
}
