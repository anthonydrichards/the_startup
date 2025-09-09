package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

const crlf = "\r\n"
const headerDelimiter = ":"

var errMalformedHeader = fmt.Errorf("malformed header")
var errWhitespaceInHeaderName = fmt.Errorf("malformed header - trailing whitespace in header name")
var errEmptyHeaderName = fmt.Errorf("malformed header - empty header name")
var errInvalidHeaderName = fmt.Errorf("malformed header - invalid header name")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	var totalBytesParsed int
	idx := bytes.Index(data[totalBytesParsed:], []byte(crlf))
	if idx == -1 {
		return totalBytesParsed, false, nil
	}
	if idx == 0 {
		totalBytesParsed += len(crlf)
		return totalBytesParsed, true, nil
	}
	header := string(data[:idx])
	headerParts := strings.SplitN(header, headerDelimiter, 2)
	if len(headerParts) != 2 {
		return 0, false, errMalformedHeader
	}
	if headerParts[0] == "" {
		return 0, false, errEmptyHeaderName
	}
	if !IsToken(headerParts[0]) {
		return 0, false, errInvalidHeaderName
	}
	if unicode.IsSpace(rune(headerParts[0][len(headerParts[0])-1])) {
		return 0, false, errWhitespaceInHeaderName
	}
	fieldName := strings.TrimSpace(headerParts[0])
	fieldValue := strings.TrimSpace(headerParts[1])
	h.Set(fieldName, fieldValue)
	totalBytesParsed += idx + len(crlf)
	return totalBytesParsed, false, nil
}

const tokenChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+-.^_`|~"

func IsToken(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !strings.ContainsRune(tokenChars, r) {
			return false
		}
	}
	return true
}

func (h Headers) Set(name, value string) {
	existingValue, ok := h[strings.ToLower(name)]
	if ok {
		h[strings.ToLower(name)] = existingValue + ", " + value
		return
	}
	h[strings.ToLower(name)] = value
}

func (h Headers) Get(name string) (string, bool) {
	value, ok := h[strings.ToLower(name)]
	return value, ok
}
