package ifce

type Unmarshaler interface {
	UnmarshalPHP([]byte) error
}
