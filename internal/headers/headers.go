package headers

import (
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

const crlf = "\r\n"
const headerDelimiter = ":"

var errMalformedHeader = fmt.Errorf("malformed header")
var errWhitespaceInHeaderName = fmt.Errorf("malformed header - whitespace in header name")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	var totalBytesParsed int
	if !strings.Contains(string(data), crlf+crlf) {
		return 0, false, nil
	}
	for _, header := range strings.Split(string(data), crlf) {
		if header == crlf {
			return totalBytesParsed + len(crlf), true, nil
		}
		headerParts := strings.SplitN(header, headerDelimiter, 2)
		if len(headerParts) != 2 {
			return 0, false, errMalformedHeader
		}
		if unicode.IsSpace(rune(headerParts[0][len(headerParts[0])-1])) {
			return 0, false, errWhitespaceInHeaderName
		}

	}
	return n, true, nil
}
