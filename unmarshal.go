package phpserialize

type Unmarshaler interface {
	UnmarshalPHP([]byte) error
}

func Unmarshal(data []byte, v interface{}) error {
	return unmarshal(data, v)
}

func UnmarshalWithOption(data []byte, v interface{}) error {
	return unmarshal(data, v)
}

func UnmarshalNoEscape(data []byte, v interface{}) error {
	return unmarshalNoEscape(data, v)
}
