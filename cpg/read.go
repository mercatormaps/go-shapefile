package cpg

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding"
)

// Read a .cpg file containing a character encoding.
func Read(r io.Reader) (encoding.Encoding, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if len(s) == 0 {
			continue
		}

		s = strings.ToUpper(s)
		switch s {
		case "ASCII":
			return encoding.Nop, nil
		case "UTF8":
			fallthrough
		case "UTF-8":
			return encoding.Nop, nil
		default:
			return encoding.Nop, nil
		}
	}
	return nil, fmt.Errorf("invalid format")
}
