package dbf

import (
	"golang.org/x/text/encoding"
)

// Option funcs can be passed to Scanner.Scan().
type Option func(*config)

// CharacterEncoding sets the encoding of character field values.
// By default, ASCII is assumed.
func CharacterEncoding(enc encoding.Encoding) Option {
	return func(c *config) {
		c.charEnc = enc
	}
}

// FilterFields allows filtering by field name.
// If this option is used, only these fields will be returned in the Record.
// Without this option, all available fields are returned.
func FilterFields(names ...string) Option {
	return func(c *config) {
		c.fields = names
	}
}

// Config for dbf parsing.
type config struct {
	charEnc encoding.Encoding
	fields  []string
}

// CharacterEncoding returns the configured encoding.
func (c *config) CharacterEncoding() encoding.Encoding {
	return c.charEnc
}

// FilteredFields returns the configured field names.
func (c *config) FilteredFields() []string {
	return c.fields
}

func defaultConfig() *config {
	return &config{
		charEnc: encoding.Nop,
	}
}
