package phpserialize

type Unmarshaler interface {
	UnmarshalPHP([]byte) error
}

func Unmarshal(data []byte, v any) error {
	return unmarshal(data, v)
}
