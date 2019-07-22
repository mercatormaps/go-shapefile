package field

import (
	"bytes"
	"strings"

	"golang.org/x/text/encoding"
)

// Character field is a string of characters.
type Character struct {
	Field
	String string
}

// DecodeCharacter decodes a single character field with the specified encoding.
func DecodeCharacter(buf []byte, name string, decoder *encoding.Decoder) (*Character, error) {
	val := bytes.Trim(buf, "\x00")

	decoded, err := decoder.Bytes(val)
	if err != nil {
		return nil, err
	}

	return &Character{
		Field:  Field{name: name},
		String: strings.TrimSpace(string(decoded)),
	}, nil
}

// Value returns the field value.
func (c *Character) Value() interface{} {
	return c.String
}
