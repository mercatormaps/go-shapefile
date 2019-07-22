package cpg_test

import (
	"bytes"
	"testing"

	"github.com/mercatormaps/go-shapefile/cpg"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding"
)

func TestRead(t *testing.T) {
	tests := []struct {
		in  string
		out encoding.Encoding
	}{
		{"ASCII", encoding.Nop},
		{"UTF-8", encoding.Nop},
		{"UTF8", encoding.Nop},
		{"iicsa", encoding.Nop},
	}

	for _, tt := range tests {
		in := "\n\t   " + tt.in + "   \t\n"
		out, err := cpg.Read(bytes.NewBuffer([]byte(in)))
		require.NoError(t, err)
		require.Equal(t, tt.out, out)
	}
}
