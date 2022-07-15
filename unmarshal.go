package phpserialize

type Unmarshaler interface {
	UnmarshalPHP([]byte) error
}

func Unmarshal(data []byte, v any) error {
	if len(data) == 0 {
		return nil
	}

	return unmarshal(data, v)
}
