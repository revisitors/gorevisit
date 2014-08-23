package gorevisit

import (
	"errors"
)

var (
	//ErrUnsupportedType is returned when a Transform does not support the type(s) passed to it
	ErrUnsupportedType = errors.New("unsupported type")
)

// Transformer interface transforms an APIMsg into another APIMsg
type Transformer interface {
	Transform(*APIMsg) (*APIMsg, error)
}
